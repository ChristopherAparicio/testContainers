package sql

import (
	"context"
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	tinyError "github.com/christapa/tinyurl/pkg/error"
	tinySql "github.com/christapa/tinyurl/pkg/sql"
)

var (
	_ domain.UrlRepository = &TinyUrlSqlRepository{}
)

// TinyUrlSqlRepository is a struct that represent the UrlRepository
// Only for SQL database (no cache)
// Implement UrlRepository interface
type TinyUrlSqlRepository struct {
	querier tinySql.Querier
}

func NewUrlSqlRepository(querier tinySql.Querier) *TinyUrlSqlRepository {
	return &TinyUrlSqlRepository{querier: querier}
}

// URL represents a URL record in the database
type URL struct {
	ShortenURL  string
	OriginalURL string
	Counter     int
	Expiration  time.Time
}

func (u *TinyUrlSqlRepository) StoreUrl(ctx context.Context, url domain.Url) (domain.Url, error) {
	result, err := u.querier.ExecContext(ctx,
		`INSERT INTO urls (shorten_url, original_url, counter, expiration_date) VALUES ($1, $2, $3, $4)`,
		url.ShortenURL,
		url.OriginalURL,
		url.Counter,
		url.Expiration,
	)
	if err != nil {
		return domain.Url{}, sqlToDomainError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.Url{}, tinyError.New(tinyError.Internal, err.Error())
	}

	if rowsAffected == 0 {
		return domain.Url{}, tinyError.New(tinyError.Internal, "failed to insert url")
	}

	return url, nil
}

func (u *TinyUrlSqlRepository) GetUrl(ctx context.Context, shortUrl string) (domain.Url, error) {
	rows, err := u.querier.QueryContext(ctx,
		"SELECT shorten_url, original_url, counter, expiration_date FROM urls WHERE shorten_url = $1 LIMIT 1",
		shortUrl)
	if err != nil {
		return domain.Url{}, sqlToDomainError(err)
	}

	defer rows.Close()

	var url domain.Url

	var found bool
	if rows.Next() {
		found = true
		err = rows.Scan(&url.ShortenURL, &url.OriginalURL, &url.Counter, &url.Expiration)
		if err != nil {
			return domain.Url{}, tinyError.New(tinyError.Internal, err.Error())
		}
	}

	if !found {
		return domain.Url{}, tinyError.New(tinyError.NotFound, "not found")
	}

	return url, nil
}

func (u *TinyUrlSqlRepository) IncrementCounter(ctx context.Context, shortUrl string) error {
	result, err := u.querier.ExecContext(ctx,
		"UPDATE urls SET counter = counter + 1 WHERE shorten_url = $1",
		shortUrl,
	)

	if err != nil {
		return sqlToDomainError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tinyError.New(tinyError.Internal, err.Error())
	}

	if rowsAffected == 0 {
		return tinyError.New(tinyError.NotFound, "not found")
	}

	return nil
}

func (u *TinyUrlSqlRepository) DeleteUrl(ctx context.Context, shortUrl string) error {
	result, err := u.querier.ExecContext(ctx,
		"DELETE FROM urls WHERE shorten_url = $1",
		shortUrl,
	)

	if err != nil {
		return sqlToDomainError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tinyError.New(tinyError.Internal, err.Error())
	}

	if rowsAffected == 0 {
		return tinyError.New(tinyError.NotFound, "not found")
	}

	return nil
}
