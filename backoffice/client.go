package backoffice

import (
	"net/http"
)

type client struct {
	httpClient      *http.Client
	authToken       string
	refreshOnExpiry bool
}

func New(opts ...Option) (Client, error) {
	c := &client{
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

type Option func(c *client) error

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) error {
		c.httpClient = httpClient
		return nil
	}
}

func WithAuthToken(authToken string) Option {
	return func(c *client) error {
		c.authToken = authToken
		return nil
	}
}
