package apps

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

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
	"github.com/hertz-contrib/limiter"
	"github.com/hertz-contrib/logger/accesslog"
	hertztracingzap "github.com/hertz-contrib/obs-opentelemetry/logging/zap"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/pprof"
	nacos_hertz "github.com/hertz-contrib/registry/nacos/v2"
	"github.com/hertz-contrib/requestid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/redis/go-redis/v9"
)

type OnStartup func(ctx SmartContext)
type OnShutdown route.CtxCallback

func NewHertzApp(appName string, genContext GenContext, options ...Option) *server.Hertz {
	_logger := logk.NewZapLogger("./logs/"+appName+".log", 20, 10, 7, hlog.LevelDebug)
	hlog.SetLogger(_logger)
	//
	opt := &option{validators: make(map[string]func(args ...interface{}) error)}
	for _, o := range options {
		o(opt)
	}
	//
	defaultPort := 8089
	conf := &Configuration{
		RateLimitEnabled: true,
		Port:             defaultPort,
		MaxRequestBodyKB: 50,
		DistributeLock:   false,
		LogLevel:         logk.Level(hlog.LevelDebug),
		Profile:          Profiling{Type: Pprof, Enabled: true},
	}
	//
	env := os.Getenv("SMART_APP_ENV")
	cfgFile := "application.yaml"
	if env != "" {
		cfgFile = fmt.Sprintf("application.%s.yaml", env)
	}
	hlog.Infof("load app configuration file: %s", cfgFile)
	_loader := loaders.NewCompositeLoader(
		loaders.NewLocalLoader("./"+cfgFile),
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
	db, err := rdbs.NewRDB(conf.RDB, opt.plugins...)
	if err != nil {
		hlog.Fatalf("failed to create rdb instance: %v", err)
		return nil
	}
	if err = db.AutoMigrate(opt.models...); err != nil {
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
	//
	for k, v := range opt.validators {
		validateConfig.MustRegValidateFunc(k, v)
	}
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
		if nnc, err = nacos.NewNamingClientWithConfig(*conf.Nacos, conf.LogLevel.String()); err != nil {
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
	h.Use(requestid.New())
	//////////////////////////////////////////////////////////////////////////////////////////
	if len(conf.AccessLog) > 0 {
		if strings.EqualFold(conf.AccessLog, "default") {
			conf.AccessLog = "[${time}] | ${requestId} | ${status} | [r:${bytesReceived},s:${bytesSent}] | - ${latency} ${method} ${contentType} ${path}"
		}
		accesslog.Tags["requestId"] = func(output accesslog.Buffer, c *app.RequestContext, data *accesslog.Data, extraParam string) (int, error) {
			return output.WriteString(requestid.Get(c))
		}
		accesslog.Tags["contentType"] = func(output accesslog.Buffer, c *app.RequestContext, data *accesslog.Data, extraParam string) (int, error) {
			return output.WriteString(string(c.ContentType()))
		}
		h.Use(accesslog.New(accesslog.WithFormat("[AccessLog] " + conf.AccessLog)))
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	if conf.RateLimitEnabled {
		h.Use(limiter.AdaptiveLimit())
	}
	//////////////////////////////////////////////////////////////////////////////////////////
	if tracerConfig != nil {
		hlog.Info("logger with tracing")
		hlog.SetLogger(hertztracingzap.NewLogger(hertztracingzap.WithLogger(_logger)))
		h.Use(hertztracing.ServerMiddleware(tracerConfig))
	}
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
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		ctx.Next(context.WithValue(c, ContextKeySmart, smartCtx))
	})
	//
	if opt.middlewares != nil {
		h.Use(opt.middlewares(smartCtx)...)
	}
	//
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
	}, route.CtxCallback(opt.shutdownHandle))
	//
	opt.startupHandle(smartCtx)
	//
	initProfile(conf, h.Group("/mgr", opt.mgrAuth...), smartCtx)
	//
	return h
}

func initProfile(conf *Configuration, g *route.RouterGroup, _ SmartContext) {
	//
	if conf.Profile.Type == None || !conf.Profile.Enabled {
		hlog.Infof("app profiling is not enabled")
	} else {
		if conf.Profile.Type == Pprof {
			pprof.RouteRegister(g)
		} else if conf.Profile.Type == FGprof {
			pprof.FgprofRouteRegister(g)
		}
	}
}
