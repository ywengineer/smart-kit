package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/jwt"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"go.uber.org/zap"
	"time"
)

type JwtConfig struct {
	Realm       string        `yaml:"realm" json:"realm"`
	Key         string        `yaml:"key" json:"key"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout"`
	MaxRefresh  time.Duration `yaml:"max_refresh" json:"max-refresh"`
	IdentityKey string        `yaml:"identity-key" json:"identity_key"`
}

func NewJwt(cfg JwtConfig) *jwt.HertzJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       cfg.Realm,
		Key:         []byte(cfg.Key),
		Timeout:     cfg.Timeout,
		MaxRefresh:  cfg.MaxRefresh,
		IdentityKey: cfg.IdentityKey,
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			return data != nil
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, pkg.ApiError(message))
		},
	})
	if err != nil {
		hlog.Fatal("create JWT middleware Error", zap.Error(err))
	}
	return authMiddleware
}
