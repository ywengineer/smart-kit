package caches

import (
	"context"
	"errors"
	"strings"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
)

type Option[T any] func(*redisCache[T])

func WithMemory[T any]() Option[T] {
	return func(c *redisCache[T]) {
		c.m = NewLocalCache[T](1 << 32)
	}
}

type redisCache[T any] struct {
	client redis.UniversalClient
	l      *singleflight.Group
	m      Cache[T]
}

func NewRedisCache[T any](client redis.UniversalClient, opts ...Option[T]) Cache[T] {
	c := &redisCache[T]{client: client, l: &singleflight.Group{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (r *redisCache[T]) GetWithLoader(key string, loader func() (T, time.Duration, error)) (T, error) {
	t, err := r.Get(key)
	if err == nil {
		return t, nil
	}
	t, err, _ = r.l.Do(key, func() (any, error) {
		v, ttl, le := loader()
		if le != nil {
			return nil, le
		}
		// put to cache
		return v, r.PutWithTtl(key, v, ttl)
	})
	return t, err
}

func (r *redisCache[T]) Get(key string) (T, error) {
	// get from memory cache
	if r.m != nil {
		t, err := r.m.Get(key)
		if err == nil {
			return t, nil
		}
	}
	// get from redis
	c, err := r.client.Get(context.Background(), key).Result()
	var value T
	if errors.Is(err, redis.Nil) {
		return value, ErrNotFound
	}
	if err := sonic.UnmarshalString(c, &value); err != nil {
		return value, err
	}
	return value, nil
}

func (r *redisCache[T]) Put(key string, value T) error {
	return r.PutWithTtl(key, value, 0)
}

func (r *redisCache[T]) PutWithTtl(key string, value T, ttl time.Duration) error {
	// put to memory cache
	if r.m != nil {
		_ = r.m.PutWithTtl(key, value, utilk.Max(0, ttl))
	}
	// put to redis
	json, err := sonic.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, json, ttl).Err()
}

func (r *redisCache[T]) Invalidate(key string) error {
	// invalidate memory cache
	if r.m != nil {
		_ = r.m.Invalidate(key)
	}
	return r.client.Del(context.Background(), key).Err()
}

func (r *redisCache[T]) InvalidatePrefix(prefix string) error {
	ctx := context.Background()
	// invalidate memory cache
	if r.m != nil {
		_ = r.m.InvalidatePrefix(prefix)
	}
	// invalidate redis cache
	pattern := lo.If(strings.HasSuffix(prefix, ":"), prefix+"*").Else(prefix + ":*") // 要删除的键前缀
	count := int64(10)                                                               // 每次迭代返回的键数量
	// 迭代扫描所有匹配的键
	cursor := uint64(0)
	for {
		// 执行 SCAN 命令
		keys, newCursor, err := r.client.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return err
		}
		// 批量删除当前批次的键
		if len(keys) > 0 {
			_, err = r.client.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}
		// 4. 游标为 0 时，迭代结束
		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	//
	return nil
}
