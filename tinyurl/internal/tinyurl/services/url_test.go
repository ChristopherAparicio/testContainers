package services

import (
	"context"
	"testing"
	"time"

	"github.com/christapa/tinyurl/internal/tinyurl/domain"
	"github.com/christapa/tinyurl/internal/tinyurl/domain/mocks"
	"github.com/stretchr/testify/assert"
)

// TODO : Increase coverage

func TestCreateShortenURLHappyPath(t *testing.T) {
	urlRepositoryMock := mocks.NewUrlRepository(t)

	url, err := domain.NewURL("https://www.google.com", time.Time{})
	if err != nil {
		t.Errorf("Error while creating a new URL: %v", err)
	}

	ctx := context.Background()

	urlRepositoryMock.On("StoreUrl", ctx, url).Return(url, nil)

	// Create a new URL service
	urlService := NewUrlService(urlRepositoryMock)

	// Call the CreateShortenUrl function
	urlCreated, err := urlService.CreateShortenUrl(ctx, url.OriginalURL, url.Expiration)
	if err != nil {
		t.Errorf("Error while creating a shorten URL: %v", err)
	}

	urlRepositoryMock.AssertExpectations(t)

	assert.Equal(t, url, urlCreated, "The two urls should be equal (Create/Created)")
}
