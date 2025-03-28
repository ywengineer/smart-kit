package apps

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
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/cors"
	nacos_hertz "github.com/hertz-contrib/registry/nacos/v2"
	"github.com/hertz-contrib/requestid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/pkg/loaders"
	"github.com/ywengineer/smart-kit/pkg/locks"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"github.com/ywengineer/smart-kit/pkg/nacos"
	"github.com/ywengineer/smart-kit/pkg/nets"
	"github.com/ywengineer/smart-kit/pkg/rdbs"
	"github.com/ywengineer/smart-kit/pkg/rpcs"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/ywengineer/smart-kit/pkg/validator"
	"net/http"
	"strings"
	"time"
)

type OnStartup func(ctx SmartContext)
type OnShutdown route.CtxCallback

func NewHertzApp(appName string,
	genContext GenContext,
	startup OnStartup,
	shutdown OnShutdown,
	rdbModels ...interface{},
) *server.Hertz {
	hlog.SetLogger(logk.NewLogger("./logs/"+appName+".log", 20, 10, 7, hlog.LevelDebug))
	//
	defaultPort := 8089
	conf := &Configuration{Port: defaultPort, MaxRequestBodyKB: 50, DistributeLock: false, LogLevel: logk.Level(hlog.LevelDebug)}
	_loader := loaders.NewLocalLoader("./application.yaml")
	if err := _loader.Load(conf); err != nil {
		hlog.Fatalf("failed to load application.yaml: %v", err)
	}
	hlog.SetLevel(hlog.Level(conf.LogLevel))
	//
	if err := _loader.Watch(context.Background(), func(data string) error {
		if err := _loader.Unmarshal([]byte(data), conf); err != nil {
			return err
		}
		hlog.SetLevel(hlog.Level(conf.LogLevel))
		return nil
	}); err != nil {
		hlog.Fatalf("failed to watch app configuration: %v", err)
	}
	//
	if (conf.RegistryInfo != nil || conf.DiscoveryEnable) && conf.Nacos == nil {
		hlog.Fatalf("enable service registry or discovery. but not found nacos configuration")
		return nil
	}
	conf.Port = utilk.Min(utilk.Max(conf.Port, 1), 65535)
	// redis
	var redisClient redis.UniversalClient
	var lockMgr locks.Manager
	if len(conf.Redis) > 0 {
		redisClient = utilk.NewRedis(conf.Redis)
	}
	if !conf.DistributeLock {
		lockMgr = locks.NewSystemLockManager()
	} else if redisClient == nil {
		hlog.Fatalf("can not create distribute lock, because of redis client is nil")
		return nil
	} else {
		lockMgr = locks.NewRedisLockManager(redislock.New(redisClient))
	}
	// rational database
	db, err := rdbs.NewRDB(conf.RDB)
	if err != nil {
		hlog.Fatalf("failed to create rdb instance: %v", err)
		return nil
	}
	if err = db.AutoMigrate(rdbModels...); err != nil {
		hlog.Fatalf("failed to start orm migrate: %v", err)
		return nil
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
			return nil
		}
	}
	if conf.RegistryInfo != nil {
		addr := conf.RegistryInfo.Addr
		if len(addr) == 0 {
			addr = nets.GetDefaultIpv4()
		}
		if strings.ContainsRune(addr, ':') == false {
			addr = fmt.Sprintf("%s:%d", addr, conf.Port)
		}
		//
		conf.RegistryInfo.Addr = addr
		//
		sOption = append(sOption, server.WithRegistry(nacos_hertz.NewNacosRegistry(nnc), &registry.Info{
			ServiceName: conf.RegistryInfo.ServiceName,
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      utilk.Max(1, conf.RegistryInfo.Weight),
			Tags:        conf.RegistryInfo.Tags,
		}))
	} else {
		conf.RegistryInfo = &ServiceInfo{
			ServiceName: "appName",
			Addr:        fmt.Sprintf("%s:%d", nets.GetDefaultIpv4(), conf.Port),
			Weight:      1,
			Tags:        map[string]string{},
		}
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
	var rpc rpcs.Rpc
	var rpcClientInfo = rpcs.RpcClientInfo{
		ClientName:     conf.RegistryInfo.String(),
		MaxRetry:       1,
		Delay:          time.Millisecond * 10,
		MaxConnPerHost: 256,
	}
	if conf.DiscoveryEnable {
		rpc, err = rpcs.NewHertzRpc(nacos_hertz.NewNacosResolver(nnc), rpcClientInfo)
	} else {
		rpc, err = rpcs.NewHertzRpc(nil, rpcClientInfo)
	}
	if err != nil {
		hlog.Fatalf("failed to create rpc: %v", err)
		return nil
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	smartCtx := genContext(
		db,
		redisClient,
		lockMgr,
		NewJwt(*conf.Jwt, nil),
		rpc,
		conf,
	)
	//
	h.Use(requestid.New())
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Next(context.WithValue(c, ContextKeySmart, smartCtx))
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
		if nnc != nil {
			nnc.CloseClient()
		}
	}, route.CtxCallback(shutdown))
	//
	startup(smartCtx)
	return h
}
