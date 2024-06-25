package services

import (
	"context"
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	"github.com/christapa/tinyurl/internal/tinyurl/usecases"
	"github.com/christapa/tinyurl/pkg/logger"
)

var (
	// Ensure UrlService implements the domain.UrlUseCase interface
	_ usecases.URL = (*UrlService)(nil)
)

type UrlService struct {
	repository domain.UrlRepository
}

func NewUrlService(repository domain.UrlRepository) *UrlService {
	return &UrlService{
		repository: repository,
	}
}

// CreateShortenUrl : Shorten the url and store it in the database
func (u *UrlService) CreateShortenUrl(ctx context.Context, url string, expiration time.Time) (domain.Url, error) {
	newUrl, err := domain.NewURL(url, expiration)
	if err != nil {
		return domain.Url{}, err
	}

	_, err = u.repository.StoreUrl(ctx, newUrl)
	if err != nil {
		return domain.Url{}, err
	}

	return newUrl, nil
}

// GetOriginalUrl : Give back the original url from the shorten url
// Each time the original url is retrieved, the counter is incremented
// If the url is expired, it is deleted from the database
func (u *UrlService) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	url, err := u.repository.GetUrl(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	err = u.assessUrl(url)
	if err != nil {
		return "", err
	}

	err = u.repository.IncrementCounter(ctx, shortUrl)
	if err != nil {
		// Functionnal : Does it needs to be strongly consistent ?
		logger.Errorf("Error incrementing counter for shortUrl %s: %v", shortUrl, err)
	}

	return url.OriginalURL, nil
}

// GetUrlMetadata : Give back the metadata of the url
func (u *UrlService) GetURLMetadata(ctx context.Context, shortUrl string) (domain.Url, error) {
	url, err := u.repository.GetUrl(ctx, shortUrl)
	if err != nil {
		return domain.Url{}, err
	}

	err = u.assessUrl(url)
	if err != nil {
		return domain.Url{}, err
	}

	return url, nil
}

func (u *UrlService) assessUrl(url domain.Url) error {
	if url.IsExpired() {
		err := u.repository.DeleteUrl(context.Background(), url.ShortenURL)
		if err != nil {
			return err
		}

		return domain.NewNotFoundError()
	}

	return nil
}
