package locks

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

type redisLockMgr struct {
	cli *redislock.Client
}

func NewRedisLockManager(cli *redislock.Client) Manager {
	return &redisLockMgr{cli: cli}
}

func (r *redisLockMgr) Obtain(ctx context.Context, key string, ttl time.Duration, opt *redislock.Options) (Lock, error) {
	return r.cli.Obtain(ctx, key, ttl, opt)
}

func (r *redisLockMgr) Close() {
}
