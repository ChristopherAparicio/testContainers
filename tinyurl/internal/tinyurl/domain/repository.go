package domain

import "context"

type UrlRepository interface {
	StoreUrl(ctx context.Context, url Url) (Url, error)
	GetUrl(ctx context.Context, shortUrl string) (Url, error)
	IncrementCounter(ctx context.Context, shortUrl string) error
	DeleteUrl(ctx context.Context, shortUrl string) error
}
