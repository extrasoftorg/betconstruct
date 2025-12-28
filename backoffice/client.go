package backoffice

import (
	"net/http"
)

type client struct {
	httpClient *http.Client
	authToken  string
}

func New(authToken string) Client {
	return &client{
		httpClient: http.DefaultClient,
		authToken:  authToken,
	}
}
