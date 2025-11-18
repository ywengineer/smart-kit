package caches

import (
	"errors"
	"time"
)

var ErrUnsupported = errors.New("unsupported operation")
var ErrNotFound = errors.New("key not found")

type Cache[T any] interface {
	Get(key string) (T, error)
	Put(key string, value T) error
	PutWithTtl(key string, value T, ttl time.Duration) error
	Invalidate(key string) error
	InvalidatePrefix(prefix string) error
}
