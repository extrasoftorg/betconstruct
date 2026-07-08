package pool

import (
	"context"
	"sync"
	"time"
)

type RateLimiter interface {
	IsLimited(ctx context.Context, token string) (bool, error)
	SetLimited(ctx context.Context, token string) error
}

const defaultRateLimiterDuration = 4 * time.Minute

var _ RateLimiter = &defaultRateLimiter{}

type defaultRateLimiter struct {
	mu       sync.Mutex
	limited  map[string]time.Time
	duration time.Duration
}

type NewDefaultRateLimiterOptions struct {
	Duration time.Duration
}

func NewDefaultRateLimiter(opts NewDefaultRateLimiterOptions) *defaultRateLimiter {
	r := &defaultRateLimiter{
		limited:  make(map[string]time.Time),
		duration: defaultRateLimiterDuration,
	}
	if opts.Duration != 0 {
		r.duration = opts.Duration
	}
	return r
}

func (r *defaultRateLimiter) IsLimited(ctx context.Context, token string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.limited[token]; ok {
		passed := time.Since(r.limited[token]) > r.duration
		if passed {
			delete(r.limited, token)
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func (r *defaultRateLimiter) SetLimited(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limited[token] = time.Now()
	return nil
}
