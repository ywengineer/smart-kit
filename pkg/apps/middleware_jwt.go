package apps

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"go.uber.org/zap"
	"time"
)

type JwtConfig struct {
	Realm       string        `yaml:"realm" json:"realm"`
	Key         string        `yaml:"key" json:"key"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout"`
	MaxRefresh  time.Duration `yaml:"max-refresh" json:"max_refresh"`
	IdentityKey string        `yaml:"identity-key" json:"identity_key"`
}

type Authentication func(data interface{}, ctx context.Context, c *app.RequestContext) bool

func NewJwt(cfg JwtConfig, auth Authentication) *jwt.HertzJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:        cfg.Realm,
		Key:          []byte(cfg.Key),
		Timeout:      cfg.Timeout,
		MaxRefresh:   cfg.MaxRefresh,
		IdentityKey:  cfg.IdentityKey,
		Authorizator: auth,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(jwt.MapClaims); ok {
				return v
			} else if v, ok := data.(map[string]interface{}); ok {
				return v
			}
			return jwt.MapClaims{}
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			hlog.Infof("auth failed: %s", message)
			c.JSON(code, ApiError("common.err.invalid_token"))
		},
	})
	if err != nil {
		hlog.Fatal("create JWT middleware Error", zap.Error(err))
	}
	return authMiddleware
}
