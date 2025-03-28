package locks

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/bsm/redislock"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"io"
	"sync"
	"time"
)

var ErrNotObtained = errors.New("failed to obtain lock")
var ErrExpired = errors.New("expired lock")

type sysLockMgr struct {
	tmpMu   sync.Mutex
	tmp     []byte
	buckets []*sysLockBucket
	cap     int
}

type sysLock struct {
	key      string
	tk       string
	meta     string
	expireAt int64
	_bucket  *sysLockBucket
}

func (s *sysLock) Key() string {
	return s.key
}

func (s *sysLock) Token() string {
	return s.tk
}

func (s *sysLock) Metadata() string {
	return s.meta
}

func (s *sysLock) TTL(ctx context.Context) (time.Duration, error) {
	return time.Duration(utilk.Max(s.expireAt-time.Now().Unix(), 0)) * time.Second, nil
}

func (s *sysLock) Refresh(ctx context.Context, ttl time.Duration, opt *redislock.Options) error {
	if dur, err := s.TTL(ctx); err != nil {
		return err
	} else if dur <= 0 {
		return ErrExpired
	} else {
		s.expireAt += ttl.Milliseconds()
		return nil
	}
}

func (s *sysLock) Release(ctx context.Context) error {
	if s._bucket != nil {
		delete(s._bucket.ch, s.key)
		s._bucket.lockPool.Put(s)
		s._bucket = nil
	}
	return nil
}

func NewSystemLockManager() Manager {
	lm := &sysLockMgr{buckets: []*sysLockBucket{}, cap: 32}
	// default 16 buckets
	for range lm.cap {
		lm.buckets = append(lm.buckets, &sysLockBucket{ch: make(map[string]*sysLock, 10), lockPool: &sync.Pool{New: func() interface{} {
			return &sysLock{}
		}}})
	}
	return lm
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

	retry := opt.RetryStrategy

	// make sure we don't retry forever
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(ttl))
		defer cancel()
	}

	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()
	//
	bucket := r.buckets[XXHash(key)&uint64(r.cap-1)]
	for {
		if l, _ := bucket.obtain(key, token, opt.Metadata, ttl); l != nil {
			return l, nil
		}
		//
		backoff := retry.NextBackoff()
		if backoff < 1 {
			return nil, ErrNotObtained
		}
		//
		if ticker == nil {
			ticker = time.NewTicker(backoff)
		} else {
			ticker.Reset(backoff)
		}
		//
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
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

type sysLockBucket struct {
	sync.Mutex
	ch       map[string]*sysLock
	lockPool *sync.Pool
}

func (r *sysLockBucket) obtain(key, token, meta string, ttlVal time.Duration) (Lock, error) {
	r.Lock()
	defer r.Unlock()
	now := time.Now().UnixMilli()
	if l, ok := r.ch[key]; !ok {
		l = r.lockPool.Get().(*sysLock)
		l._bucket = r
		l.tk, l.meta, l.expireAt, l.key = token, meta, now+ttlVal.Milliseconds(), key
		r.ch[key] = l
		return l, nil
	} else if l.Token() == token && l.Metadata() == meta && l.expireAt > now {
		l.expireAt += ttlVal.Milliseconds()
		return l, nil
	} else {
		return nil, ErrNotObtained
	}
}
