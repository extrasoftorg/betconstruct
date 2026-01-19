package accounts

import "net/http"

type client struct {
	httpClient *http.Client
}

func New(opts ...Option) Client {
	c := &client{
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(c *client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}
