package pool

import (
	"context"
	"sync"
	"time"
)

type Pool interface {
	// GetAuthToken returns a non-rate-limited auth token from the pool.
	// If no auth token is available, it returns nil.
	// It uses a round-robin algorithm to get the next available token.
	GetAuthToken(ctx context.Context) *AuthToken
	SetRateLimited(ctx context.Context, token string) error
}

type pool struct {
	mu         sync.Mutex
	authTokens []AuthToken
	nextIndex  int
}

func (p *pool) GetAuthToken(ctx context.Context) *AuthToken {
	p.mu.Lock()
	defer p.mu.Unlock()

	n := len(p.authTokens)
	for i := range n {
		idx := (p.nextIndex + i) % n
		if !p.authTokens[idx].IsRateLimited() {
			p.authTokens[idx].setLastUsed(time.Now())
			p.nextIndex = (idx + 1) % n
			return &p.authTokens[idx]
		}
	}
	return nil
}

func (p *pool) SetRateLimited(ctx context.Context, token string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.authTokens {
		if p.authTokens[i].String() == token {
			p.authTokens[i].setRateLimited(time.Now())
			break
		}
	}
	return nil
}

func New(authTokens ...AuthToken) Pool {
	return &pool{
		authTokens: authTokens,
	}
}
