package session

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// Implement the encoding.BinaryUnmarshaler interface
type Session struct {
	SessionID string
	Username  string
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *Session) getKey() string {
	return s.SessionID
}

func (s *Session) getValue() *Session {
	return s
}

type RedisSessionRepository struct {
	client *redis.Client
	TTL    time.Duration
}

func NewRedisSessionRepository(client *redis.Client, ttl time.Duration) *RedisSessionRepository {
	return &RedisSessionRepository{
		client: client,
		TTL:    ttl,
	}
}

func (r *RedisSessionRepository) Save(ctx context.Context, session *Session) error {
	return r.client.Set(ctx, session.getKey(), session.getValue(), r.TTL).Err()
}

func (r *RedisSessionRepository) Get(ctx context.Context, sessionID string) (*Session, error) {
	session := &Session{}
	err := r.client.Get(ctx, sessionID).Scan(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
