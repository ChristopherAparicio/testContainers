package user

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	postgresmodules "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func initPostgresContainer(ctx context.Context, dbConfig DBConfig) (*postgresmodules.PostgresContainer, error) {
	postgresContainer, err := postgresmodules.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres"),
		postgresmodules.WithInitScripts(filepath.Join("db.sql")),
		postgresmodules.WithDatabase(dbConfig.Database),
		postgresmodules.WithUsername(dbConfig.User),
		postgresmodules.WithPassword(dbConfig.Password),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Cmd: []string{"-c", "log_statement=all"},
			},
		}),

		/*
			testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
				Reuse:            true,
				ContainerRequest: testcontainers.ContainerRequest{
					// Name: "postgres-container",
				},
			}),
		*/

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	return postgresContainer, err
}

func getPostgresContainerPort(ctx context.Context, postgresContainer testcontainers.Container) (int, error) {
	natPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get container mapped port: %s", err)
	}

	return natPort.Int(), nil
}

func newPostgresConnectionToTestContainerWithConnString(ctx context.Context) (*sql.DB, error) {
	dbConfig := DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Database: "example",
	}

	postgresContainer, err := initPostgresContainer(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	conn, err := NewConn(connectionString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newPostgresConnectionToTestContainerWithDBConfig(ctx context.Context, dbConfig DBConfig) (*sql.DB, error) {
	postgresContainer, err := initPostgresContainer(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	postgresPort, err := getPostgresContainerPort(ctx, postgresContainer)
	if err != nil {
		return nil, err
	}

	dbConfig.Port = postgresPort
	conn, err := NewConnWithDbConfig(dbConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func freezeContainer() {
	if flag.Lookup("test.timeout").Value.String() == "0s" {
		cancelChan := make(chan os.Signal, 1)

		signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

		_ = <-cancelChan
		return
	}
}

func TestPostgreslSessionRepository_Scenario(t *testing.T) {
	ctx := context.Background()
	defer freezeContainer()

	databaseConnection, err := newPostgresConnectionToTestContainerWithConnString(ctx)

	// databaseConnection, err := newPostgresConnectionToTestContainerWithDBConfig(ctx, DBConfig{
	// 	Host:     "localhost",
	// 	User:     "postgres",
	// 	Password: "postgres",
	// 	Database: "example",
	// })
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	t.Run("CRD User", func(t *testing.T) {
		userRepository := NewSqlUserRepository(databaseConnection)

		user := &User{
			Email:        "test@test.com",
			PasswordHash: "password",
			LastLogin:    time.Now().UTC(),
		}

		_, err := userRepository.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		savedUser, err := userRepository.GetUserByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		assert.Equal(t, user.Email, savedUser.Email)
		assert.Equal(t, user.PasswordHash, savedUser.PasswordHash)
		assert.Equal(t, user.LastLogin.Unix(), savedUser.LastLogin.Unix())

		err = userRepository.DeleteUserByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		savedUser, err = userRepository.GetUserByEmail(ctx, user.Email)
		assert.NotNil(t, err, "error should not be nil")
	})

}

func freezeContainerHelper() {

}
