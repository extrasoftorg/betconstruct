package crm

import (
	"context"
	"net/http"
)

type client struct {
	httpClient        *http.Client
	authToken         string
	betconstructToken string
	refreshOnExpiry   bool
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := &client{
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.betconstructToken != "" {
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Option func(c *client) error

func WithAuthToken(authToken string) Option {
	return func(c *client) error {
		c.authToken = authToken
		return nil
	}
}

func WithBetconstructToken(betconstructToken string) Option {
	return func(c *client) error {
		c.betconstructToken = betconstructToken
		return nil
	}
}

func WithRefreshOnExpiry() Option {
	return func(c *client) error {
		c.refreshOnExpiry = true
		return nil
	}
}
