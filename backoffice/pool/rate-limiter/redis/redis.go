package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/extrasoftorg/betconstruct/backoffice/pool"
	"github.com/redis/go-redis/v9"
)

const (
	defaultKey      = "auth_token_pool:ratelimit"
	defaultDuration = 4 * time.Minute
)

var _ pool.RateLimiter = &rateLimiter{}

type rateLimiter struct {
	client   *redis.Client
	key      string
	duration time.Duration
}

type Option func(*rateLimiter)

func New(ctx context.Context, redis *redis.Client, opts ...Option) (*rateLimiter, error) {
	if redis == nil {
		return nil, ErrRedisRequired
	}
	if redis.Ping(ctx).Err() != nil {
		return nil, ErrRedisNotAvailable
	}

	r := &rateLimiter{
		client:   redis,
		key:      defaultKey,
		duration: defaultDuration,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

func WithKey(key string) Option {
	return func(r *rateLimiter) {
		r.key = key
	}
}

func WithDuration(duration time.Duration) Option {
	return func(r *rateLimiter) {
		r.duration = duration
	}
}

func (r *rateLimiter) IsLimited(ctx context.Context, token string) (bool, error) {
	val, err := r.client.Get(ctx, r.keyForToken(token)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return val != "", nil
}

func (r *rateLimiter) SetLimited(ctx context.Context, token string) error {
	return r.client.Set(ctx, r.keyForToken(token), time.Now().Unix(), r.duration).Err()
}

func (r *rateLimiter) keyForToken(token string) string {
	return fmt.Sprintf("%s:%s", r.key, token)
}
