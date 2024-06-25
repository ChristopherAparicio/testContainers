package domain

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type MockShortenURLGenerator struct{}

func (m MockShortenURLGenerator) GenerateShortenURL(url string) string {
	return url
}

func TestNewURL(t *testing.T) {
	type args struct {
		originalURL string
		expiration  time.Time
	}
	tests := []struct {
		name string
		args args
		want Url
		err  error
	}{
		{
			name: "Test NewURL",
			args: args{
				originalURL: "https://www.google.com",
			},
			want: Url{
				ShortenURL:  "https://www.google.com",
				OriginalURL: "https://www.google.com",
				Counter:     0,
			},
			err: nil,
		},
		{
			name: "Expired URL",
			args: args{
				originalURL: "https://www.google.com",
				expiration:  time.Now().Add(-time.Hour),
			},
			want: Url{},
			err:  NewInvalidInputError("expiration date is in the past"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetShortenURLGenerator(MockShortenURLGenerator{})

			got, err := NewURL(tt.args.originalURL, tt.args.expiration)
			if err != nil && !errors.Is(err, tt.err) {
				t.Errorf("NewURL() = %v, want %v", err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultShortenURLGenerator(t *testing.T) {
	url := "https://www.google.com"
	want := "rGu2aeQO"
	got := defaultShortenURLGenerator{}.GenerateShortenURL(url)
	if got != want {
		t.Errorf("GenerateShortenURL() = %v, want %v", got, want)
	}
}

func TestUrl_IsExpired(t *testing.T) {
	type fields struct {
		ShortenURL  string
		OriginalURL string
		Counter     int
		Expiration  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Not expired",
			fields: fields{
				Expiration: time.Now().Add(time.Hour),
			},
			want: false,
		},
		{
			name: "Expired",
			fields: fields{
				Expiration: time.Now().Add(-time.Hour),
			},
			want: true,
		},
		{
			name: "No expiration",
			fields: fields{
				Expiration: time.Time{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			U := &Url{
				ShortenURL:  tt.fields.ShortenURL,
				OriginalURL: tt.fields.OriginalURL,
				Counter:     tt.fields.Counter,
				Expiration:  tt.fields.Expiration,
			}
			if got := U.IsExpired(); got != tt.want {
				t.Errorf("Url.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrl_IncrementCounter(t *testing.T) {
	type fields struct {
		ShortenURL  string
		OriginalURL string
		Counter     int
		Expiration  time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Increment counter",
			fields: fields{
				Counter: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Url{
				ShortenURL:  tt.fields.ShortenURL,
				OriginalURL: tt.fields.OriginalURL,
				Counter:     tt.fields.Counter,
				Expiration:  tt.fields.Expiration,
			}
			u.IncrementCounter()

			if u.Counter != 1 {
				t.Errorf("Counter = %v, want %v", u.Counter, 1)
			}
		})
	}
}
