package backoffice

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type stubResponse struct {
	status int
	body   string
}

// stubTransport serves canned responses in order and records the token used for
// every request. Once the canned responses run out, the last one repeats.
type stubTransport struct {
	mu        sync.Mutex
	responses []stubResponse
	tokens    []string
}

func (t *stubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tokens = append(t.tokens, req.Header.Get("Authentication"))

	idx := min(len(t.tokens)-1, len(t.responses)-1)
	r := t.responses[idx]

	return &http.Response{
		StatusCode: r.status,
		Body:       io.NopCloser(strings.NewReader(r.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func (t *stubTransport) calls() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]string(nil), t.tokens...)
}

// stubTokenSource rotates over tokens and skips the ones marked as limited.
type stubTokenSource struct {
	mu      sync.Mutex
	tokens  []string
	next    int
	limited map[string]bool
}

func (s *stubTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tokens {
		idx := (s.next + i) % len(s.tokens)
		if s.limited[s.tokens[idx]] {
			continue
		}
		s.next = (idx + 1) % len(s.tokens)
		return s.tokens[idx], nil
	}
	return "", ErrRateLimited
}

func (s *stubTokenSource) MarkLimited(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.limited == nil {
		s.limited = make(map[string]bool)
	}
	s.limited[token] = true
	return nil
}

func newTestClient(t *testing.T, tr http.RoundTripper, opts ...Option) Client {
	t.Helper()

	c, err := New(append([]Option{WithHTTPClient(&http.Client{Transport: tr})}, opts...)...)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return c
}

const playerBody = `{"Data":{"Id":42,"Login":"john"},"HasError":false}`

// Non-2xx responses must surface as sentinel errors, not as JSON decode errors.
func TestMakeRequest_StatusCodeMapping(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr error
	}{
		{"bad request", http.StatusBadRequest, ErrBadRequest},
		{"not found", http.StatusNotFound, ErrNotFound},
		{"method not allowed", http.StatusMethodNotAllowed, ErrMethodNotAllowed},
		{"internal server error", http.StatusInternalServerError, ErrInternalServerError},
		{"bad gateway", http.StatusBadGateway, ErrBadGateway},
		{"service unavailable", http.StatusServiceUnavailable, ErrServiceUnavailable},
		{"unexpected status", http.StatusTeapot, ErrUnexpectedStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &stubTransport{responses: []stubResponse{{status: tt.status, body: "<html>error</html>"}}}
			c := newTestClient(t, tr, WithAuthToken("token-1"))

			if _, err := c.GetPlayer(context.Background(), 42); !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if calls := tr.calls(); len(calls) != 1 {
				t.Fatalf("made %d requests, want 1 (no retry on non-token errors)", len(calls))
			}
		})
	}
}

func TestMakeRequest_Success(t *testing.T) {
	tr := &stubTransport{responses: []stubResponse{{status: http.StatusOK, body: playerBody}}}
	c := newTestClient(t, tr, WithAuthToken("token-1"))

	player, err := c.GetPlayer(context.Background(), 42)
	if err != nil {
		t.Fatalf("failed to get player: %v", err)
	}
	if player.ID != 42 || player.Username != "john" {
		t.Fatalf("got player %+v, want ID 42 and username john", player)
	}
}

// A rejected token must be marked limited and the request retried with the next one.
func TestMakeRequest_RetriesWithNextToken(t *testing.T) {
	tr := &stubTransport{responses: []stubResponse{
		{status: http.StatusForbidden, body: ""},
		{status: http.StatusOK, body: playerBody},
	}}
	ts := &stubTokenSource{tokens: []string{"token-1", "token-2"}}
	c := newTestClient(t, tr, WithTokenSource(ts))

	if _, err := c.GetPlayer(context.Background(), 42); err != nil {
		t.Fatalf("failed to get player: %v", err)
	}

	calls := tr.calls()
	if len(calls) != 2 {
		t.Fatalf("made %d requests, want 2", len(calls))
	}
	if calls[0] == calls[1] {
		t.Fatalf("retried with the same token %q, want a different one", calls[0])
	}
	if !ts.limited[calls[0]] {
		t.Fatalf("token %q was not marked as limited", calls[0])
	}
}

// With a single token there is nothing to rotate to, so the real cause must be
// returned after exactly one request.
func TestMakeRequest_StaticTokenDoesNotRetry(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr error
	}{
		{"unauthorized", http.StatusUnauthorized, ErrUnauthorized},
		{"forbidden", http.StatusForbidden, ErrForbidden},
		{"too many requests", http.StatusTooManyRequests, ErrTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &stubTransport{responses: []stubResponse{{status: tt.status, body: ""}}}
			c := newTestClient(t, tr, WithAuthToken("token-1"))

			if _, err := c.GetPlayer(context.Background(), 42); !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v, want %v", err, tt.wantErr)
			}
			if calls := tr.calls(); len(calls) != 1 {
				t.Fatalf("made %d requests, want 1", len(calls))
			}
		})
	}
}

// Running out of attempts must report what the server actually said.
func TestMakeRequest_ExhaustedAttemptsReportLastStatus(t *testing.T) {
	tr := &stubTransport{responses: []stubResponse{{status: http.StatusTooManyRequests, body: ""}}}
	ts := &stubTokenSource{tokens: []string{"token-1", "token-2", "token-3"}}
	c := newTestClient(t, tr, WithTokenSource(ts), WithMaxAttempts(2))

	if _, err := c.GetPlayer(context.Background(), 42); !errors.Is(err, ErrTooManyRequests) {
		t.Fatalf("got error %v, want %v", err, ErrTooManyRequests)
	}
	if calls := tr.calls(); len(calls) != 2 {
		t.Fatalf("made %d requests, want 2", len(calls))
	}
}

// When every token is limited the token source decides, and its error surfaces.
func TestMakeRequest_AllTokensLimited(t *testing.T) {
	tr := &stubTransport{responses: []stubResponse{{status: http.StatusForbidden, body: ""}}}
	ts := &stubTokenSource{tokens: []string{"token-1", "token-2"}}
	c := newTestClient(t, tr, WithTokenSource(ts))

	if _, err := c.GetPlayer(context.Background(), 42); !errors.Is(err, ErrRateLimited) {
		t.Fatalf("got error %v, want %v", err, ErrRateLimited)
	}
	if calls := tr.calls(); len(calls) != 2 {
		t.Fatalf("made %d requests, want 2", len(calls))
	}
}

func TestStatusError(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		if err := statusError(status); err != nil {
			t.Fatalf("statusError(%d) = %v, want nil", status, err)
		}
	}
	if err := statusError(http.StatusNotFound); !errors.Is(err, ErrNotFound) {
		t.Fatalf("statusError(404) = %v, want %v", err, ErrNotFound)
	}
}

func TestNew_RequiresTokenSource(t *testing.T) {
	if _, err := New(); !errors.Is(err, ErrMissingTokenSource) {
		t.Fatalf("got error %v, want %v", err, ErrMissingTokenSource)
	}
}
