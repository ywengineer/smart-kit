package apps

import (
	"context"
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/loaders"
	"gitee.com/ywengineer/smart-kit/pkg/locks"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/nets"
	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"gitee.com/ywengineer/smart-kit/pkg/validator"
	"github.com/bsm/redislock"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/binding"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/logger/accesslog"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/pprof"
	nacos_hertz "github.com/hertz-contrib/registry/nacos/v2"
	"github.com/hertz-contrib/requestid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/redis/go-redis/v9"
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
	hlog.SetLogger(logk.NewZapLogger("./logs/"+appName+".log", 20, 10, 7, hlog.LevelDebug))
	//
	defaultPort := 8089
	conf := &Configuration{Port: defaultPort, MaxRequestBodyKB: 50, DistributeLock: false, LogLevel: logk.Level(hlog.LevelDebug), Profile: Profiling{Type: Pprof, Enabled: true, AuthDownload: true, Prefix: "/mgr/prof"}}
	_loader := loaders.NewCompositeLoader(
		loaders.NewLocalLoader("./application.yaml"),
		loaders.NewEnvLoader(),
	)
	if err := _loader.Load(conf); err != nil {
		hlog.Fatalf("failed to load application.yaml: %v", err)
	}
	hlog.SetLevel(hlog.Level(conf.LogLevel))
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
		server.WithTraceLevel(stats.Level(conf.TraceLevel)),
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	var nnc naming_client.INamingClient
	if conf.Nacos != nil {
		conf.Nacos.Cluster = utilk.DefaultIfEmpty(conf.Nacos.Cluster, "DEFAULT")
		conf.Nacos.Group = utilk.DefaultIfEmpty(conf.Nacos.Group, "DEFAULT_GROUP")
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
		sOption = append(sOption, server.WithRegistry(nacos_hertz.NewNacosRegistry(nnc, nacos_hertz.WithRegistryCluster(conf.Nacos.Cluster), nacos_hertz.WithRegistryGroup(conf.Nacos.Group)), &registry.Info{
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
	var tracerConfig *hertztracing.Config
	if stats.Level(conf.TraceLevel) != stats.LevelDisabled {
		var tracer config.Option
		tracer, tracerConfig = hertztracing.NewServerTracer()
		sOption = append(sOption, tracer)
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	h := server.Default(sOption...)
	//////////////////////////////////////////////////////////////////////////////////////////
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
	if tracerConfig != nil {
		h.Use(hertztracing.ServerMiddleware(tracerConfig))
	}
	if len(conf.AccessLog) > 0 {
		if strings.EqualFold(conf.AccessLog, "default") {
			conf.AccessLog = "[${time}] | ${requestId} | ${status} | [r:${bytesReceived},s:${bytesSent}] | - ${latency} ${method} ${path}"
		}
		accesslog.Tags["requestId"] = func(output accesslog.Buffer, c *app.RequestContext, data *accesslog.Data, extraParam string) (int, error) {
			return output.WriteString(requestid.Get(c))
		}
		h.Use(accesslog.New(accesslog.WithFormat("[AccessLog] " + conf.AccessLog)))
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
		rpc, err = rpcs.NewHertzRpc(nacos_hertz.NewNacosResolver(nnc, nacos_hertz.WithResolverCluster(conf.Nacos.Cluster), nacos_hertz.WithResolverGroup(conf.Nacos.Group)), rpcClientInfo)
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
	//
	initProfile(conf, h, smartCtx)
	//
	return h
}

func initProfile(conf *Configuration, h *server.Hertz, ctx SmartContext) {
	//
	if conf.Profile.Type == None || !conf.Profile.Enabled {
		hlog.Infof("app profiling is not enabled")
	} else {
		if len(conf.Profile.Prefix) == 0 {
			hlog.Infof("app profile path is not set, default is /mgr/prof")
			conf.Profile.Prefix = "/mgr/prof"
		}
		var g *route.RouterGroup
		if conf.Profile.AuthDownload {
			g = h.Group(conf.Profile.Prefix, ctx.TokenInterceptor())
		} else {
			g = h.Group(conf.Profile.Prefix)
		}
		if conf.Profile.Type == Pprof {
			pprof.RouteRegister(g)
		} else if conf.Profile.Type == FGprof {
			pprof.FgprofRouteRegister(g)
		}
	}
}
