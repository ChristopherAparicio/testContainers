package usecases

import (
	"context"
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
)

type URL interface {
	CreateShortenUrl(ctx context.Context, url string, expiration time.Time) (domain.Url, error)
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
	GetURLMetadata(ctx context.Context, url string) (domain.Url, error)
}
