package caches

import (
	"github.com/dgraph-io/ristretto/v2"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"sync"
)

var cache *ristretto.Cache[string, []byte]
var s sync.Once

func init() {
	s.Do(func() {
		var err error
		cache, err = ristretto.NewCache(&ristretto.Config[string, []byte]{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})
		if err != nil {
			panic(err)
		} else {
			go func() {
				defer cache.Close()
				<-utilk.WatchQuitSignal()
			}()
		}
	})
}

func Get(key string) ([]byte, bool) {
	return cache.Get(key)
}

func Put(key string, value []byte, cost int64) bool {
	return cache.Set(key, value, cost)
}
