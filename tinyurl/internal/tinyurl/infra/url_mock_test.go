package sql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	tinyError "github.com/christapa/tinyurl/pkg/error"
	"github.com/stretchr/testify/assert"
)

var databaseConnection *sql.DB

func TestMockStoreUrlHappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	url := domain.Url{
		ShortenURL:  "rGu2aeQO",
		OriginalURL: "https://www.google.com",
		Counter:     0,
		Expiration:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs(url.ShortenURL, url.OriginalURL, url.Counter, url.Expiration).
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := NewUrlSqlRepository(db)

	_, err = service.StoreUrl(context.Background(), url)
	if err != nil {
		t.Fatalf("Failed to store url : %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMockStoreUrlShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	url := domain.Url{
		ShortenURL:  "rGu2aeQO",
		OriginalURL: "https://www.google.com",
		Counter:     0,
		Expiration:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs(url.ShortenURL, url.OriginalURL, url.Counter, url.Expiration).
		WillReturnError(errors.New("some error"))

	service := NewUrlSqlRepository(db)

	_, err = service.StoreUrl(context.Background(), url)

	assert.ErrorIs(t, err, tinyError.New(tinyError.Internal, "some error"))
}

func TestMockStoreUrlShouldReturnEmptyRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	url := domain.Url{
		ShortenURL:  "rGu2aeQO",
		OriginalURL: "https://www.google.com",
		Counter:     0,
		Expiration:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO urls").
		WithArgs(url.ShortenURL, url.OriginalURL, url.Counter, url.Expiration).
		WillReturnResult(sqlmock.NewResult(1, 0))

	service := NewUrlSqlRepository(db)

	_, err = service.StoreUrl(context.Background(), url)

	assert.ErrorIs(t, err, tinyError.New(tinyError.Internal, "failed to insert url"))
}

// TestSqlScenario is a test that will test the SQL scenario
// It will create a new url, get it, increment the counter and delete it
func TestSqlDatabaseScenario(t *testing.T) {
	if databaseConnection == nil {
		t.Fatalf("Database connection is nil")
	}

	tx, err := databaseConnection.BeginTx(context.Background(), nil)
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
}

func TestDeleteUrlShouldReturnNotFound(t *testing.T) {
	if databaseConnection == nil {
		t.Fatalf("Database connection is nil")
	}

	tx, err := databaseConnection.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to start transaction : %v", err)
	}

	defer tx.Commit()

	service := NewUrlSqlRepository(tx)

	err = service.DeleteUrl(context.Background(), "unknown")
	assert.Error(t, err)

	assert.True(t, errors.Is(err, tinyError.New(tinyError.NotFound, "not found")))
}
