package domain

import (
	"time"
)

type Url struct {
	ShortenURL  string
	OriginalURL string
	Counter     int
	Expiration  time.Time
}

func NewURL(originalURL string, expiration time.Time) (Url, error) {
	if err := isValidOriginalURL(originalURL); err != nil {
		return Url{}, NewInvalidInputError(err.Error())
	}

	if !expiration.IsZero() && expiration.Before(time.Now()) {
		return Url{}, NewInvalidInputError("expiration date is in the past")
	}

	shortenURL := ShortenURLGeneratorInstance.GenerateShortenURL(originalURL)

	return Url{
		ShortenURL:  string(shortenURL[:]),
		OriginalURL: originalURL,
		Counter:     0,
		Expiration:  expiration,
	}, nil
}

func (u *Url) IncrementCounter() {
	u.Counter++
}

func (u *Url) IsExpired() bool {

	return !u.Expiration.IsZero() && u.Expiration.Before(time.Now())
}
