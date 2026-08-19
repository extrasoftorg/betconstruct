package backoffice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://backofficewebadmin.betconstruct.com/api/en"
)

type response[T any] struct {
	Data         T      `json:"Data"`
	HasError     bool   `json:"HasError"`
	AlertMessage string `json:"AlertMessage"`
}

func (c *client) attempt(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	token string,
) (*http.Response, error) {
	fullURL := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func makeRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body []byte,
	c *client,
) (*T, error) {
	var (
		lastToken  string
		lastStatus int
	)

	for attempt := range c.maxAttempts {
		token, err := c.ts.Token(ctx)
		if err != nil {
			return nil, err
		}

		// If the token is the same as the last one, we can assume that the token source is not providing a new token and we should not retry.
		if attempt > 0 && token == lastToken {
			return nil, statusError(lastStatus)
		}
		lastToken = token

		resp, err := c.attempt(ctx, method, path, bytes.NewReader(body), token)
		if err != nil {
			return nil, err
		}

		if isTokenRejection(resp.StatusCode) {
			lastStatus = resp.StatusCode
			drainAndClose(resp.Body)
			_ = c.ts.MarkLimited(ctx, token)
			continue
		}

		if err := statusError(resp.StatusCode); err != nil {
			drainAndClose(resp.Body)
			return nil, err
		}

		return decode[T](resp)
	}

	return nil, statusError(lastStatus)
}

func decode[T any](resp *http.Response) (*T, error) {
	defer drainAndClose(resp.Body)

	var data response[T]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	} else if data.HasError {
		return nil, errors.New(data.AlertMessage)
	}

	return &data.Data, nil
}

// statusError maps an HTTP status code to a sentinel error. It returns nil for
// any 2xx status.
func statusError(statusCode int) error {
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return nil
	}

	switch statusCode {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusMethodNotAllowed:
		return ErrMethodNotAllowed
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		return ErrInternalServerError
	case http.StatusBadGateway:
		return ErrBadGateway
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return ErrUnexpectedStatus
	}
}

func isTokenRejection(statusCode int) bool {
	return statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden || statusCode == http.StatusTooManyRequests
}

func drainAndClose(body io.ReadCloser) {
	io.Copy(io.Discard, io.LimitReader(body, 64<<10))
	body.Close()
}
