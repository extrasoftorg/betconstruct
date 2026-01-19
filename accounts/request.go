package accounts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://api.accounts-bc.com/"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrInternalServerError = errors.New("internal server error")
	ErrBadGateway          = errors.New("bad gateway")
	ErrServiceUnavailable  = errors.New("service unavailable")
	ErrUnexpectedStatus    = errors.New("unexpected status")
)

func makeRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	c *client,
	marshal func(r io.Reader) error,
) (*T, error) {
	fullURL := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/html")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, ErrBadRequest
		case http.StatusUnauthorized:
			return nil, ErrUnauthorized
		case http.StatusForbidden:
			return nil, ErrForbidden
		case http.StatusNotFound:
			return nil, ErrNotFound
		case http.StatusMethodNotAllowed:
			return nil, ErrMethodNotAllowed
		case http.StatusTooManyRequests:
			return nil, ErrTooManyRequests
		case http.StatusInternalServerError:
			return nil, ErrInternalServerError
		case http.StatusBadGateway:
			return nil, ErrBadGateway
		case http.StatusServiceUnavailable:
			return nil, ErrServiceUnavailable
		default:
			return nil, ErrUnexpectedStatus
		}
	}

	if marshal != nil {
		err := marshal(resp.Body)
		return nil, err
	}

	return nil, nil
}
