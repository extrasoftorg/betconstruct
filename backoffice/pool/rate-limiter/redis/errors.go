package redis

import "errors"

var (
	ErrRedisRequired     = errors.New("redis is required")
	ErrRedisNotAvailable = errors.New("redis is not available")
)
