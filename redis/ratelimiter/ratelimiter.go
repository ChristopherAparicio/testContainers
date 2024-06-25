package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client   *redis.Client
	duration time.Duration
	rate     int64
}

func NewRateLimiter(client *redis.Client, windowsDuration time.Duration, rate int64) *RateLimiter {
	return &RateLimiter{
		client:   client,
		duration: windowsDuration,
		rate:     rate,
	}
}

func (r *RateLimiter) RateLimiter(ctx context.Context, ip string) (bool, error) {

	key := fmt.Sprintf("%s", ip)

	// Increment the counter for this IP in the current minute window
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set the expiration for this key if this is the first request in the current minute
	if count == 1 {
		r.client.Expire(ctx, key, r.duration)
	}

	// Check if the request count exceeds the limit
	if count > int64(r.rate) {
		return true, nil
	}

	return false, nil

}
