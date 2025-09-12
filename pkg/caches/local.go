package caches

import (
	"sync"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/dgraph-io/ristretto/v2"
)

var cache *ristretto.Cache[string, []byte]
var s sync.Once

func init() {
	s.Do(func() {
		var err error
		cache, err = NewCache[[]byte](1 << 30)
		if err != nil {
			panic(err)
		}
	})
}

func NewCache[T any](capacity int64) (*ristretto.Cache[string, T], error) {
	c, err := ristretto.NewCache(&ristretto.Config[string, T]{
		NumCounters: 1e7,      // number of keys to track frequency of (10M).
		MaxCost:     capacity, // maximum cost of cache (1GB).
		BufferItems: 64,       // size of Get buffers.
	})
	if err != nil {
		return nil, err
	} else {
		go func() {
			defer c.Close()
			<-utilk.WatchQuitSignal()
		}()
	}
	return c, nil
}

func Get(key string) ([]byte, bool) {
	return cache.Get(key)
}

func Put(key string, value []byte, cost int64) bool {
	return cache.Set(key, value, cost)
}
