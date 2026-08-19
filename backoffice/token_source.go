package backoffice

import "context"

type TokenSource interface {
	Token(ctx context.Context) (string, error)
	MarkLimited(ctx context.Context, token string) error
}

type staticTokenSource struct {
	token string
}

func (s *staticTokenSource) Token(ctx context.Context) (string, error) {
	return s.token, nil
}

func (s *staticTokenSource) MarkLimited(ctx context.Context, token string) error {
	return nil
}

var _ TokenSource = &staticTokenSource{}
