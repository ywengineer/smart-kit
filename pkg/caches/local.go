package caches

import (
	"errors"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/dgraph-io/ristretto/v2"
)

type localCache[T any] struct {
	c *ristretto.Cache[string, T]
}

func NewLocalCache[T any](capacity int64) Cache[T] {
	c, err := ristretto.NewCache(&ristretto.Config[string, T]{
		NumCounters: 1e7,      // number of keys to track frequency of (10M).
		MaxCost:     capacity, // maximum cost of cache (1GB).
		BufferItems: 64,       // size of Get buffers.
	})
	if err != nil {
		logk.Fatalf("local cache init fail: %v", err)
		return nil
	} else {
		go func() {
			defer c.Close()
			<-utilk.WatchQuitSignal()
		}()
	}
	return &localCache[T]{c: c}
}

func (l *localCache[T]) Get(key string) (T, error) {
	t, ok := l.c.Get(key)
	if !ok {
		return t, ErrNotFound
	}
	return t, nil
}

func (l *localCache[T]) Put(key string, value T) error {
	if l.c.Set(key, value, 1) == false {
		return errors.New("failed to put key in local cache")
	}
	return nil
}

func (l *localCache[T]) PutWithTtl(key string, value T, ttl time.Duration) error {
	if l.c.SetWithTTL(key, value, 1, ttl) == false {
		return errors.New("failed to put key with ttl in local cache")
	}
	return nil
}

func (l *localCache[T]) Invalidate(key string) error {
	l.c.Del(key)
	return nil
}

func (l *localCache[T]) InvalidatePrefix(prefix string) error {
	l.c.Clear()
	return nil
}
