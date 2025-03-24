// Code generated by hertz generator.

package main

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/cors"
	hertzzap "github.com/hertz-contrib/logger/zap"
	nacos_hertz "github.com/hertz-contrib/registry/nacos/v2"
	"github.com/hertz-contrib/requestid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/redis/go-redis/v9"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"github.com/ywengineer/smart-kit/passport/pkg/lock"
	"github.com/ywengineer/smart-kit/passport/pkg/middleware"
	"github.com/ywengineer/smart-kit/passport/pkg/validator"
	"github.com/ywengineer/smart-kit/pkg/nacos"
	"github.com/ywengineer/smart/loader"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	hlog.SetLogger(logger)
	hlog.SetLevel(hlog.LevelDebug)
	//
	defaultPort := 8089
	conf := &Configuration{Port: defaultPort, MaxRequestBodyKB: 50, DistributeLock: false, LogLevel: zap.DebugLevel}
	_loader := loader.NewLocalLoader("./application.yaml")
	if err := _loader.Load(conf); err != nil {
		hlog.Fatalf("failed to load application.yaml: %v", err)
	}
	hlog.SetLevel(hlog.Level(conf.LogLevel + 1))
	//
	if err := _loader.Watch(context.Background(), func(data string) error {
		if err := _loader.Unmarshal([]byte(data), conf); err != nil {
			return err
		}
		hlog.SetLevel(hlog.Level(conf.LogLevel + 1))
		return nil
	}); err != nil {
		hlog.Fatalf("failed to watch app configuration: %v", err)
	}
	//
	if (conf.RegistryInfo != nil || conf.DiscoveryEnable) && conf.Nacos == nil {
		hlog.Fatalf("enable service registry or discovery. but not found nacos configuration")
		return
	}
	conf.Port = utility.MinInt(utility.MaxInt(conf.Port, 1), 65535)
	// redis
	var redisClient redis.UniversalClient
	var lockMgr lock.Manager
	if len(conf.Redis) > 0 {
		redisClient = utility.NewRedis(conf.Redis)
	}
	if !conf.DistributeLock {
		lockMgr = lock.NewSystemLockManager()
	} else if redisClient == nil {
		panic("can not create distribute lock, because of redis client is nil")
	} else {
		lockMgr = lock.NewRedisLockManager(redislock.New(redisClient))
	}
	// rational database
	db, err := utility.NewRDB(conf.RDB)
	if err != nil {
		hlog.Fatalf("failed to create rdb instance: %v", err)
	}
	err = db.AutoMigrate(
		&model2.Passport{},
		&model2.PassportPunish{},
		&model2.PassportBinding{},
		&model2.WhiteList{},
		&model2.MgrUser{},
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
	sOption := []config.Option{
		server.WithBindConfig(bindConfig),
		server.WithHostPorts(fmt.Sprintf(":%d", conf.Port)),
		server.WithBasePath(conf.BasePath),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(conf.MaxRequestBodyKB * 1024), // KB
		server.WithValidateConfig(validateConfig),
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	var nnc naming_client.INamingClient
	if conf.Nacos != nil {
		if nnc, err = nacos.NewNacosNamingClient(conf.Nacos.Ip, conf.Nacos.Port, conf.Nacos.ContextPath, conf.Nacos.TimeoutMs, conf.Nacos.Namespace, conf.Nacos.User, conf.Nacos.Password, conf.LogLevel.String()); err != nil {
			hlog.Fatalf("failed to create nacos client: %v", err)
			return
		}
	}
	if conf.RegistryInfo != nil {
		sOption = append(sOption, server.WithRegistry(nacos_hertz.NewNacosRegistry(nnc), &registry.Info{
			ServiceName: conf.RegistryInfo.ServiceName,
			Addr:        utils.NewNetAddr("tcp", ""),
			Weight:      utility.MaxInt(1, conf.RegistryInfo.Weight),
			Tags:        conf.RegistryInfo.Tags,
		}))
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	h := server.Default(sOption...)
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
		lockMgr,
		middleware.NewJwt(*conf.Jwt, nil),
	)
	//
	sqlRunner(db)
	//
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

func sqlRunner(db *gorm.DB) {
	utility.DefaultLogger().Info(db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.WithContext(context.Background()).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "passport"}},                            // 冲突字段（唯一索引）
			DoUpdates: clause.AssignmentColumns([]string{"updated_at", "deleted_at"}), // 更新字段
		}).Create(&model2.WhiteList{Passport: 100})
	}))
}
