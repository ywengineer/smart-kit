package caches

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type RedisCache[T any] struct {
	client redis.UniversalClient
}

func NewRedisCache[T any](client redis.UniversalClient) Cache[T] {
	return &RedisCache[T]{client: client}
}

func (r *RedisCache[T]) Get(key string) (T, error) {
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

func (r *RedisCache[T]) Put(key string, value T) error {
	return r.PutWithTtl(key, value, -1)
}

func (r *RedisCache[T]) PutWithTtl(key string, value T, ttl time.Duration) error {
	json, err := sonic.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), key, json, ttl).Err()
}

func (r *RedisCache[T]) Invalidate(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache[T]) InvalidatePrefix(prefix string) error {
	ctx := context.Background()
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
