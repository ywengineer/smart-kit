package pkg

import (
	"github.com/bsm/redislock"
	"github.com/hertz-contrib/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SmartContext interface {
	Rdb() *gorm.DB
	Redis() redis.UniversalClient
	DistributeLock() *redislock.Client
	Jwt() *jwt.HertzJWTMiddleware
}

type defaultContext struct {
	rdb     *gorm.DB
	redis   redis.UniversalClient
	redLock *redislock.Client
	_jwt    *jwt.HertzJWTMiddleware
}

func NewDefaultContext(rdb *gorm.DB, redis redis.UniversalClient, redLock *redislock.Client, jwt *jwt.HertzJWTMiddleware) SmartContext {
	return &defaultContext{
		rdb:     rdb,
		redis:   redis,
		redLock: redLock,
		_jwt:    jwt,
	}
}

func (d *defaultContext) Rdb() *gorm.DB {
	return d.rdb
}

func (d *defaultContext) Redis() redis.UniversalClient {
	return d.redis
}

func (d *defaultContext) DistributeLock() *redislock.Client {
	return d.redLock
}

func (d *defaultContext) Jwt() *jwt.HertzJWTMiddleware {
	return d._jwt
}
