package pool

import "time"

const AuthTokenRateLimitDuration = 4 * time.Minute

type AuthToken struct {
	token         string
	lastUsed      *time.Time
	rateLimitedAt *time.Time
}

func (a AuthToken) String() string {
	return a.token
}

func (a AuthToken) IsRateLimited() bool {
	if a.rateLimitedAt == nil {
		return false
	}

	return time.Since(*a.rateLimitedAt) < AuthTokenRateLimitDuration
}

func (a *AuthToken) setLastUsed(lastUsed time.Time) {
	a.lastUsed = &lastUsed
}

func (a *AuthToken) setRateLimited(rateLimitedAt time.Time) {
	a.rateLimitedAt = &rateLimitedAt
}

func NewAuthToken(token string) AuthToken {
	return AuthToken{
		token: token,
	}
}
