package session

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	redisModules "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
}

var redisClient *redis.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	redisContainer, err := initCustomRedisContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to start redis container: %v", err)
	}
	redisPort, err := getRedisContainerPort(ctx, redisContainer)
	if err != nil {
		log.Fatalf("Failed to get redis container port: %v", err)
	}

	redisAddr := fmt.Sprintf("localhost:%d", redisPort)

	redisClient = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    redisAddr,
		// ClientName: "",
	})

	exitVal := m.Run()

	os.Exit(exitVal)
}

func initCustomRedisContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:         "redis-session-container",
		Image:        "docker.io/redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{Container: container}, nil
}

func initDefaultRedisContainer(ctx context.Context) (testcontainers.Container, error) {
	redisContainer, err := redisModules.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redisModules.WithLogLevel(redisModules.LogLevelDebug),
	)
	if err != nil {
		return nil, err
	}

	return redisContainer, err
}

func getRedisContainerPort(ctx context.Context, redisContainer testcontainers.Container) (int, error) {
	natPort, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("failed to get container mapped port: %s", err)
	}

	return natPort.Int(), nil
}

func TestRedisSessionRepository_Scenario(t *testing.T) {
	ctx := context.Background()

	t.Run("ReadSessionBeforeTTL", func(t *testing.T) {
		sessionRedisRepository := NewRedisSessionRepository(redisClient, 10*time.Minute)

		session := &Session{
			SessionID: uuid.New().String(),
			Username:  "test",
		}

		err := sessionRedisRepository.Save(ctx, session)
		if err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		savedSession, err := sessionRedisRepository.Get(ctx, session.SessionID)
		if err != nil {
			t.Fatalf("Failed to get session: %v", err)
		}

		assert.Equal(t, session, savedSession)
	})

	t.Run("ReadSessionAfterTTL", func(t *testing.T) {
		sessionRedisRepository := NewRedisSessionRepository(redisClient, time.Millisecond*100)

		session := &Session{
			SessionID: uuid.New().String(),
			Username:  "test",
		}

		err := sessionRedisRepository.Save(ctx, session)
		if err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		time.Sleep(time.Millisecond * 200)

		savedSession, err := sessionRedisRepository.Get(ctx, session.SessionID)
		assert.ErrorIs(t, err, redis.Nil)
		assert.Nil(t, savedSession)
	})

	time.Sleep(300 * time.Second)

}
