package backoffice

import (
	"net/http"
	"time"
)

type client struct {
	httpClient   *http.Client
	ts           TokenSource
	maxAttempts  int
	timeLocation *time.Location
}

func New(opts ...Option) (Client, error) {
	c := &client{
		httpClient:   &http.Client{},
		maxAttempts:  3,
		timeLocation: time.UTC,
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.ts == nil {
		return nil, ErrMissingTokenSource
	}

	return c, nil
}

type Option func(c *client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func WithAuthToken(authToken string) Option {
	return func(c *client) {
		c.ts = &staticTokenSource{token: authToken}
	}
}

func WithTokenSource(ts TokenSource) Option {
	return func(c *client) {
		c.ts = ts
	}
}

func WithMaxAttempts(maxAttempts int) Option {
	return func(c *client) {
		c.maxAttempts = maxAttempts
	}
}

func WithTimeLocation(location *time.Location) Option {
	return func(c *client) {
		if location == nil {
			return
		}
		c.timeLocation = location
	}
}
