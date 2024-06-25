package domain

import (
	"crypto/sha256"
	"encoding/base64"
)

var ShortenURLGeneratorInstance ShortenURLGenerator = defaultShortenURLGenerator{}

func SetShortenURLGenerator(generator ShortenURLGenerator) {
	ShortenURLGeneratorInstance = generator
}

type defaultShortenURLGenerator struct{}

func (d defaultShortenURLGenerator) GenerateShortenURL(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	hashBytes := hash.Sum(nil)

	base64Hash := base64.URLEncoding.EncodeToString(hashBytes)

	shortenedURL := base64Hash[:8]

	return shortenedURL
}

type ShortenURLGenerator interface {
	GenerateShortenURL(url string) string
}
