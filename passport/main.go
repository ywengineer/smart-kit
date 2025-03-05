// Code generated by hertz generator.

package main

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/cors"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/hertz-contrib/requestid"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"github.com/ywengineer/smart-kit/passport/pkg/middleware"
	"github.com/ywengineer/smart-kit/passport/pkg/model"
	"github.com/ywengineer/smart-kit/passport/pkg/validator"
	"github.com/ywengineer/smart/loader"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

func main() {
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/passport.out",
		MaxSize:    20,   // 一个文件最大可达 20M。
		MaxBackups: 5,    // 最多同时保存 5 个文件。
		MaxAge:     10,   // 一个文件最多可以保存 10 天。
		Compress:   true, // 用 gzip 压缩。
	}
	//
	logger := hertzzap.NewLogger(
		hertzzap.WithCoreWs(zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))),
		hertzzap.WithZapOptions(
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	)
	logger.SetLevel(hlog.LevelDebug)
	hlog.SetLogger(logger)
	//
	defaultPort := 8089
	conf := &Configuration{Port: defaultPort, MaxRequestBodyKB: 50, RedisLock: true}
	if err := loader.NewLocalLoader("./application.yaml").Load(conf); err != nil {
		hlog.Fatalf("failed to load application.yaml: %v", err)
	}
	conf.Port = utility.MinInt(utility.MaxInt(conf.Port, 1), 65535)
	// redis
	var redisClient redis.UniversalClient
	var redisLock *redislock.Client
	if len(conf.Redis) > 0 {
		redisClient = utility.NewRedis(conf.Redis)
		//
		if conf.RedisLock {
			redisLock = redislock.New(redisClient)
		}
	}
	// rational database
	db, err := utility.NewRDB(conf.RDB)
	if err != nil {
		hlog.Fatalf("failed to create rdb instance: %v", err)
	}
	err = db.AutoMigrate(
		&model.Passport{},
		&model.PassportPunish{},
		&model.PassportBinding{},
		&model.WhiteList{},
	)
	if err != nil {
		hlog.Fatalf("failed to start orm migrate: %v", err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	bindConfig := binding.NewBindConfig()
	// 默认 false，当前 Hertz Engine 下生效，多份 engine 实例之间不会冲突
	bindConfig.LooseZeroMode = true
	//////////////////////////////////////////////////////////////////////////////////////////
	validateConfig := binding.NewValidateConfig()
	validateConfig.MustRegValidateFunc("every", validator.Every)
	//////////////////////////////////////////////////////////////////////////////////////////
	h := server.Default(
		server.WithBindConfig(bindConfig),
		server.WithHostPorts(fmt.Sprintf(":%d", conf.Port)),
		server.WithBasePath(conf.BasePath),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(conf.MaxRequestBodyKB*1024), // KB
		server.WithValidateConfig(validateConfig),
	)
	if _cors := conf.Cors; _cors != nil {
		h.Use(cors.New(cors.Config{
			AllowOrigins:     _cors.AllowOrigins,
			AllowMethods:     _cors.AllowMethods,
			AllowHeaders:     _cors.AllowHeaders,
			AllowCredentials: _cors.AllowCredentials,
			ExposeHeaders:    _cors.ExposeHeaders,
			MaxAge:           _cors.MaxAge,
			AllowWildcard:    _cors.AllowWildcard,
		}))
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	smartCtx := pkg.NewDefaultContext(
		db,
		redisClient,
		redisLock,
		middleware.NewJwt(*conf.Jwt, nil),
	)
	h.Use(requestid.New())
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Next(context.WithValue(c, pkg.ContextKeySmart, smartCtx))
	})
	h.NoRoute(func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	h.NoMethod(func(c context.Context, ctx *app.RequestContext) {
		ctx.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	})
	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) {
		hlog.Info("release resource on shutdown")
		_ = redisClient.Close()
	})
	//
	register(h)
	h.Spin()
}
