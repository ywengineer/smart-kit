package app

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/pkg/lock"
	"github.com/ywengineer/smart-kit/pkg/oauths"
	"github.com/ywengineer/smart-kit/pkg/rpcs"
	"github.com/ywengineer/smart-kit/pkg/signs"
	"gorm.io/gorm"
	"strconv"
)

type SmartContext interface {
	Rdb() *gorm.DB
	Redis() redis.UniversalClient
	LockMgr() lock.Manager
	Jwt() *jwt.HertzJWTMiddleware
	TokenInterceptor() app.HandlerFunc
	GetDeviceLockKey(deviceId string) string
	GetPassportLockKey(passportId uint) string
	Rpc() rpcs.Rpc
	GetAuth(authKey string) (oauths.AuthFacade, error)
	VerifySignature(data map[string]string, sign []byte) bool
}

type defaultContext struct {
	rdb     *gorm.DB
	redis   redis.UniversalClient
	lm      lock.Manager
	_jwt    *jwt.HertzJWTMiddleware
	jwtMw   app.HandlerFunc
	mClient rpcs.Rpc
	conf    *Configuration
}

func (d *defaultContext) VerifySignature(data map[string]string, sign []byte) bool {
	return signs.VerifySignature(data, sign, d.conf.SignKey)
}

func (d *defaultContext) GetAuth(authKey string) (oauths.AuthFacade, error) {
	return d.conf.OAuth.Get(authKey)
}

func (d *defaultContext) Rpc() rpcs.Rpc {
	return d.mClient
}

func (d *defaultContext) GetDeviceLockKey(deviceId string) string {
	return "lock:device:" + deviceId
}

func (d *defaultContext) GetPassportLockKey(passportId uint) string {
	return "lock:passport:" + strconv.FormatUint(uint64(passportId), 10)
}

func NewDefaultContext(rdb *gorm.DB, redis redis.UniversalClient, lm lock.Manager, jwt *jwt.HertzJWTMiddleware, rpcClient rpcs.Rpc, conf *Configuration) SmartContext {
	return &defaultContext{
		rdb:     rdb,
		redis:   redis,
		lm:      lm,
		_jwt:    jwt,
		jwtMw:   jwt.MiddlewareFunc(),
		mClient: rpcClient,
		conf:    conf,
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

func (d *defaultContext) LockMgr() lock.Manager {
	return d.lm
}

func (d *defaultContext) Jwt() *jwt.HertzJWTMiddleware {
	return d._jwt
}
