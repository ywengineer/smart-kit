package pkg

import (
	"github.com/bsm/redislock"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

type SmartContext interface {
	Rdb() *gorm.DB
	Redis() redis.UniversalClient
	DistributeLock() *redislock.Client
	Jwt() *jwt.HertzJWTMiddleware
	TokenInterceptor() app.HandlerFunc
	GetDeviceLockKey(deviceId string) string
	GetPassportLockKey(passportId uint) string
}

type defaultContext struct {
	rdb     *gorm.DB
	redis   redis.UniversalClient
	redLock *redislock.Client
	_jwt    *jwt.HertzJWTMiddleware
	jwtMw   app.HandlerFunc
}

func (d *defaultContext) GetDeviceLockKey(deviceId string) string {
	return "lock:device:" + deviceId
}

func (d *defaultContext) GetPassportLockKey(passportId uint) string {
	return "lock:passport:" + strconv.FormatUint(uint64(passportId), 10)
}

func NewDefaultContext(rdb *gorm.DB, redis redis.UniversalClient, redLock *redislock.Client, jwt *jwt.HertzJWTMiddleware) SmartContext {
	return &defaultContext{
		rdb:     rdb,
		redis:   redis,
		redLock: redLock,
		_jwt:    jwt,
		jwtMw:   jwt.MiddlewareFunc(),
	}
}

func (d *defaultContext) TokenInterceptor() app.HandlerFunc {
	return d.jwtMw
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
