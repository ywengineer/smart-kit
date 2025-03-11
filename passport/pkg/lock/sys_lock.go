package lock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/bsm/redislock"
	"github.com/dgraph-io/ristretto/v2"
	"io"
	"strconv"
	"sync"
	"time"
)

type sysLockMgr struct {
	tmpMu sync.Mutex
	tmp   []byte
	ch    *ristretto.Cache[string, string]
}

func NewSystemLockManager() Manager {
	cache, _ := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		ShouldUpdate: func(cur, prev string) bool {
			return cur == prev
		},
	})
	return &sysLockMgr{ch: cache}
}

func (r *sysLockMgr) Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (Lock, error) {
	token := opt.Token

	// Create a random token
	if token == "" {
		var err error
		if token, err = r.randomToken(); err != nil {
			return nil, err
		}
	}

	value := token + opt.Metadata
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	retry := opt.RetryStrategy

	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(ttl))
		defer cancel()
	}

	var ticker *time.Ticker
	for {
		ok, err := r.obtain(ctx, key, value, len(token), ttlVal)
		if err != nil {
			return nil, err
		} else if ok {
			return nil, nil // &Lock{Client: c, key: key, value: value, tokenLen: len(token)}, nil
		}

		backoff := retry.NextBackoff()
		if backoff < 1 {
			return nil, redislock.ErrNotObtained
		}

		if ticker == nil {
			ticker = time.NewTicker(backoff)
			defer ticker.Stop()
		} else {
			ticker.Reset(backoff)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (r *sysLockMgr) obtain(ctx context.Context, key, value string, tokenLen int, ttlVal string) (bool, error) {
	r.tmpMu.Lock()
	defer r.tmpMu.Unlock()

	return true, nil
}

func (r *sysLockMgr) randomToken() (string, error) {
	r.tmpMu.Lock()
	defer r.tmpMu.Unlock()

	if len(r.tmp) == 0 {
		r.tmp = make([]byte, 16)
	}

	if _, err := io.ReadFull(rand.Reader, r.tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(r.tmp), nil
}

func (r *sysLockMgr) Close() {
}
