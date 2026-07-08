package pool

import (
	"context"
	"sync"
)

type Pool struct {
	mu          sync.Mutex
	authTokens  []string
	nextIndex   int
	RateLimiter RateLimiter
}

func (p *Pool) NextAuthToken(ctx context.Context) *string {
	p.mu.Lock()
	tokens := p.authTokens
	start := p.nextIndex
	p.mu.Unlock()

	n := len(tokens)
	for i := range n {
		idx := (start + i) % n
		token := tokens[idx]
		limited, err := p.RateLimiter.IsLimited(ctx, token)
		if err != nil {
			continue
		}
		if !limited {
			p.mu.Lock()
			p.nextIndex = (idx + 1) % n
			p.mu.Unlock()
			return &token
		}
	}
	return nil
}

type Option func(*Pool)

func New(authTokens []string, opts ...Option) *Pool {
	p := &Pool{
		authTokens:  authTokens,
		RateLimiter: NewDefaultRateLimiter(NewDefaultRateLimiterOptions{}),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithRateLimiter(rateLimiter RateLimiter) Option {
	return func(p *Pool) {
		p.RateLimiter = rateLimiter
	}
}
