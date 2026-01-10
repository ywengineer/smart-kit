package apps

import (
	"context"
	"strconv"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/locks"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/oauths"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/signs"
	types "gitee.com/ywengineer/smart-kit/pkg/types"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gookit/goutil/envutil"
	"github.com/hertz-contrib/jwt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _env types.Env

func init() {
	_env = types.Env(envutil.Getenv("APP_ENV", types.Production.String()))
}

func RunningIn(env types.Env) bool {
	return env == _env
}

func Env() types.Env {
	return _env
}

type GenContext func(rdb *gorm.DB, redis redis.UniversalClient, lm locks.Manager, jwt *jwt.HertzJWTMiddleware, rpcClient rpcs.Rpc, conf *Configuration) SmartContext

type SmartContext interface {
	Rdb() *gorm.DB
	Redis() redis.UniversalClient
	LockMgr() locks.Manager
	JwtIdentityKey() string
	CreateJwtToken(data interface{}) (string, time.Time, error)
	TokenInterceptor() app.HandlerFunc
	GetDeviceLockKey(deviceId string) string
	GetPassportLockKey(passportId uint) string
	Rpc() rpcs.Rpc
	GetAuth(authKey string) (oauths.AuthFacade, error)
	VerifySignature(data map[string]string, sign []byte) bool
	GetNacosConfig() nacos.Nacos
}

func GetContext(c context.Context) SmartContext {
	if r := c.Value(ContextKeySmart); r == nil {
		return nil
	} else {
		return r.(SmartContext)
	}
}

type defaultContext struct {
	rdb     *gorm.DB
	redis   redis.UniversalClient
	lm      locks.Manager
	_jwt    *jwt.HertzJWTMiddleware
	jwtMw   app.HandlerFunc
	mClient rpcs.Rpc
	conf    *Configuration
}

func (d *defaultContext) GetNacosConfig() nacos.Nacos {
	return *d.conf.Nacos
}

func (d *defaultContext) JwtIdentityKey() string {
	return d._jwt.IdentityKey
}

func (d *defaultContext) CreateJwtToken(data interface{}) (string, time.Time, error) {
	return d._jwt.TokenGenerator(data)
}

func (d *defaultContext) VerifySignature(data map[string]string, sign []byte) bool {
	if len(d.conf.SignKey) <= 0 { // 没有配置签名key，默认不验证
		return true
	}
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

func NewDefaultContext(rdb *gorm.DB, redis redis.UniversalClient, lm locks.Manager, jwt *jwt.HertzJWTMiddleware, rpcClient rpcs.Rpc, conf *Configuration) SmartContext {
	dc := &defaultContext{
		rdb:   rdb,
		redis: redis,
		lm:    lm,
		_jwt:  jwt,
		jwtMw: func(c context.Context, ctx *app.RequestContext) {
			ctx.Next(c)
		},
		mClient: rpcClient,
		conf:    conf,
	}
	if jwt != nil {
		dc.jwtMw = jwt.MiddlewareFunc()
	}
	return dc
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

func (d *defaultContext) LockMgr() locks.Manager {
	return d.lm
}
