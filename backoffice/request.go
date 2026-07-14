package backoffice

import (
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
	AlertType    string `json:"AlertType"`
}

func makeRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	c *client,
) (*T, error) {
	resp, err := doRequest[T](ctx, method, path, body, c)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func doRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	c *client,
) (*response[T], error) {
	fullURL := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	var authToken string
	if c.pool != nil {
		at := c.pool.NextAuthToken(ctx)
		if at == nil {
			return nil, ErrPoolExhausted
		}
		authToken = *at
	} else {
		authToken = c.authToken
	}
	req.Header.Set("Authentication", authToken)

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
			if c.pool != nil {
				if err := c.pool.RateLimiter.SetLimited(ctx, authToken); err != nil {
					return nil, err
				}
			}
			return nil, ErrUnauthorized
		case http.StatusForbidden:
			if c.pool != nil {
				if err := c.pool.RateLimiter.SetLimited(ctx, authToken); err != nil {
					return nil, err
				}
			}
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

	var data response[T]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	} else if data.HasError {
		return nil, errors.New(data.AlertMessage)
	}

	return &data, nil
}
