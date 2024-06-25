package sql

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	tinySql "github.com/christapa/tinyurl/pkg/sql"
	"github.com/stretchr/testify/assert"
	testcontainers "github.com/testcontainers/testcontainers-go"
	postgresmodules "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

/*
func TestMain(m *testing.M) {
	config := tinySql.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "tinyurl",
	}

	conn, err := tinySql.NewConn(config)
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	databaseConnection = conn
	defer databaseConnection.Close()

	code := m.Run()
	os.Exit(code)
}
*/

func initPostgresContainer(ctx context.Context, dbConfig tinySql.DBConfig) (testcontainers.Container, error) {
	postgresContainer, err := postgresmodules.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres"),
		postgresmodules.WithInitScripts(filepath.Join("../../../infra/sql/init.sql")),
		postgresmodules.WithDatabase(dbConfig.Database),
		postgresmodules.WithUsername(dbConfig.User),
		postgresmodules.WithPassword(dbConfig.Password),
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

func TestContainerDatabaseCreate(t *testing.T) {
	ctx := context.Background()

	dbConfig := tinySql.NewDefaultDbConfig()
	postgresContainer, err := initPostgresContainer(ctx, dbConfig)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	port, err := getPostgresContainerPort(ctx, postgresContainer)
	if err != nil {
		log.Fatalf("failed to get container mapped port: %s", err)
	}

	dbConfig.Port = port

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Test scenario
	connection, err := tinySql.NewConn(dbConfig)
	if err != nil {
		t.Fatalf("Failed to connect to database : %v", err)
	}

	tx, err := connection.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to start transaction : %v", err)
	}

	defer tx.Rollback()

	service := NewUrlSqlRepository(tx)

	url := domain.Url{
		ShortenURL:  "rGu2aeQO",
		OriginalURL: "https://www.google.com",
		Counter:     0,
		Expiration:  time.Now(),
	}

	urlCreated, err := service.StoreUrl(context.Background(), url)
	if err != nil {
		t.Fatalf("Failed to store url : %v", err)
	}

	assert.Equal(t, url, urlCreated, "The two urls should be equal (Create/Created)")

	urlGet, err := service.GetUrl(context.Background(), url.ShortenURL)
	if err != nil {
		t.Fatalf("Failed to get url : %v", err)
	}

	assert.Equal(t, url.ShortenURL, urlGet.ShortenURL, "The two shorten url should be equal (Create/Read)")
	assert.Equal(t, url.OriginalURL, urlGet.OriginalURL, "The two original urls should be equal (Create/Read)")
	assert.Equal(t, url.Counter, urlGet.Counter, "The two counters should be equal (Create/Read)")

	err = service.IncrementCounter(context.Background(), url.ShortenURL)
	if err != nil {
		t.Fatalf("Failed to increment counter : %v", err)
	}

	urlGet, err = service.GetUrl(context.Background(), url.ShortenURL)
	if err != nil {
		t.Fatalf("Failed to get url : %v", err)
	}

	assert.Equal(t, url.Counter+1, urlGet.Counter, "The url should be incremented")

	err = service.DeleteUrl(context.Background(), url.ShortenURL)
	if err != nil {
		t.Fatalf("Failed to delete url : %v", err)
	}

	fmt.Println("End test")
	time.Sleep(30 * time.Second)
}
