package locks

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/cespare/xxhash/v2"
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

func XXHash(data string) uint64 {
	// 创建 xxHash64 对象
	h := xxhash.New()
	// 写入数据
	_, _ = h.Write([]byte(data))
	// 计算哈希值
	return h.Sum64()
}

const MaximumPowerOfTwo32 = 1 << 30

// NextPowerOfTwo 函数用于计算大于等于 num 的最小 2 的 N 次幂
func NextPowerOfTwo(num uint64) uint64 {
	if num <= 0 {
		return 1
	}
	num--
	num |= num >> 1
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16
	num |= num >> 32
	if num >= MaximumPowerOfTwo32 {
		return MaximumPowerOfTwo32
	} else {
		return num + 1
	}
}
