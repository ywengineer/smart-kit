package lock

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

type Manager interface {
	Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (Lock, error)
	Close()
}

type Lock interface {
	// Key returns the redis key used by the lock.
	Key() string
	// Token returns the token value set by the lock.
	Token() string
	// Metadata returns the metadata of the lock.
	Metadata() string
	// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
	TTL(ctx context.Context) (time.Duration, error)
	// Refresh extends the lock with a new TTL.
	// May return ErrNotObtained if refresh is unsuccessful.
	Refresh(ctx context.Context, ttl time.Duration, opt *redislock.Options) error
	// Release manually releases the lock.
	// May return ErrLockNotHeld.
	Release(ctx context.Context) error
}
