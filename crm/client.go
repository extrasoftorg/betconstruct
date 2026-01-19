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
		opt(c)
	}
	if c.betconstructToken != "" {
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Option func(c *client)

func WithAuthToken(authToken string) Option {
	return func(c *client) {
		c.authToken = authToken
	}
}

func WithBetconstructToken(betconstructToken string) Option {
	return func(c *client) {
		c.betconstructToken = betconstructToken
	}
}

func WithRefreshOnExpiry() Option {
	return func(c *client) {
		c.refreshOnExpiry = true
	}
}
