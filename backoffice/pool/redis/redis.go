package redis

import (
	"context"
	"fmt"
	"time"

	poolPkg "github.com/extrasoftorg/betconstruct/backoffice/pool"
	"github.com/redis/go-redis/v9"
)

type pool struct {
	rdb          *redis.Client
	key          string
	rateLimitKey string
}

func (p *pool) GetAuthToken(ctx context.Context) *poolPkg.AuthToken {
	n, _ := p.rdb.LLen(ctx, p.key).Result()
	for range n {
		token, err := p.rdb.LMove(ctx, p.key, p.key, "RIGHT", "LEFT").Result()
		if err != nil {
			return nil
		}

		rateLimitKey := fmt.Sprintf("%s:%s", p.rateLimitKey, token)
		if p.rdb.Exists(ctx, rateLimitKey).Val() == 0 {
			at := poolPkg.NewAuthToken(token)
			return &at
		}
	}
	return nil
}

func (p *pool) SetRateLimited(ctx context.Context, token string) error {
	rateLimitKey := fmt.Sprintf("%s:%s", p.rateLimitKey, token)
	return p.rdb.Set(ctx, rateLimitKey, time.Now().Unix(), poolPkg.AuthTokenRateLimitDuration).Err()
}

func New(ctx context.Context, client *redis.Client, key string) (poolPkg.Pool, error) {
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	if key == "" {
		key = "auth_token_pool"
	}
	return &pool{
		rdb:          client,
		key:          key,
		rateLimitKey: fmt.Sprintf("%s:ratelimit", key),
	}, nil
}
