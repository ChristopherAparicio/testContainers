package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	redisModules "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
}

func initCustomValkeyContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		// Name:         "redis-session-container",
		Image:        "docker.io/valkey/valkey:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		// Reuse:            true,
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{Container: container}, nil
}

func initCustomRedisContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		// Name:         "redis-session-container",
		Image:        "docker.io/redis:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		// Reuse:            true,
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

func newRedisClientContainer(ctx context.Context, container testcontainers.Container) (*redis.Client, error) {
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	return client, nil
}

func newDefaultRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return client, nil
}

func TestRateLimiterTestContainer(t *testing.T) {
	ctx := context.Background()

	// client, err := newDefaultRedisClient()
	// assert.NoError(t, err)

	container, err := initCustomValkeyContainer(ctx)
	// container, err := initCustomRedisContainer(ctx)
	assert.NoError(t, err)

	client, err := newRedisClientContainer(ctx, container)
	assert.NoError(t, err)

	rate := int64(3)
	limiter := NewRateLimiter(client, 5*time.Second, rate)

	t.Run("No limit reached", func(t *testing.T) {
		ip := "192.168.1.50"
		hitLimit, err := limiter.RateLimiter(ctx, ip)
		assert.NoError(t, err)

		// Rate should not be exceeded
		assert.False(t, hitLimit)

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "1")
	})

	t.Run("Limit reached", func(t *testing.T) {
		ip := "192.168.1.51"
		for range rate {
			hitLimit, err := limiter.RateLimiter(ctx, ip)
			assert.NoError(t, err)
			assert.False(t, hitLimit)
		}

		hitLimit, err := limiter.RateLimiter(ctx, ip)
		assert.NoError(t, err)

		// Rate should be exceeded
		assert.True(t, hitLimit)

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "4")
	})
}
