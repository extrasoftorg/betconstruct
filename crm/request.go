package crm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL = "https://crm-t.betconstruct.com/api/en"
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

type response[T any] struct {
	Data         T      `json:"Data"`
	HasError     bool   `json:"HasError"`
	AlertMessage string `json:"AlertMessage"`
}

type makeRequestOptions struct {
	httpClient *http.Client
	authToken  string
}

func makeRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	c *client,
	marshal func(r io.Reader) error,
) (*T, error) {
	return doRequest[T](ctx, method, path, body, c, marshal, false)
}

var isRefreshing bool

func doRequest[T any](
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	c *client,
	marshal func(r io.Reader) error,
	isRetry bool,
) (*T, error) {
	fullURL := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusUnauthorized && !isRetry && c.refreshOnExpiry && !isRefreshing {
			isRefreshing = true
			err := c.Login(ctx)
			isRefreshing = false
			if err != nil {
				return nil, ErrUnauthorized
			}
			if seeker, ok := body.(io.Seeker); ok {
				seeker.Seek(0, io.SeekStart)
			}
			return doRequest[T](ctx, method, path, body, c, marshal, true)
		}

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

	var data response[T]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	} else if data.HasError {
		return nil, errors.New(data.AlertMessage)
	}

	return &data.Data, nil
}
