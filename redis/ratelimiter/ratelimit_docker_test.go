package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()

	client, err := newRedisClient()
	assert.NoError(t, err)

	rate := int64(3)
	limiter := NewRateLimiter(client, 5*time.Second, rate)

	t.Run("No limit reached", func(t *testing.T) {
		ip := "192.168.1.50"
		hitLimit, err := limiter.RateLimiter(ctx, ip)
		assert.NoError(t, err)

		// Rate should not be exceeded
		assert.False(t, hitLimit)

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "1")
	})

	t.Run("Limit reached", func(t *testing.T) {
		ip := "192.168.1.51"
		for range rate {
			hitLimit, err := limiter.RateLimiter(ctx, ip)
			assert.NoError(t, err)
			assert.False(t, hitLimit)
		}

		hitLimit, err := limiter.RateLimiter(ctx, ip)
		assert.NoError(t, err)

		// Rate should be exceeded
		assert.True(t, hitLimit)

		// Check key exists
		assert.Equal(t, client.Get(ctx, ip).Val(), "4")
	})
}
