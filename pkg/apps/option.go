package apps

import (
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type option struct {
	plugins        []gorm.Plugin
	models         []interface{}
	mgrAuth        []app.HandlerFunc
	startupHandle  OnStartup
	shutdownHandle OnShutdown
	middlewares    func(sc SmartContext) []app.HandlerFunc
	validators     map[string]func(args ...interface{}) error
	beforeMigrate  func(sc *gorm.DB) error
}

type Option func(*option)

func WithBeforeMigrate(fun func(sc *gorm.DB) error) Option {
	return func(o *option) {
		o.beforeMigrate = fun
	}
}

func WithMiddlewares(fun func(sc SmartContext) []app.HandlerFunc) Option {
	return func(o *option) {
		o.middlewares = fun
	}
}

func WithStartupHandle(handle OnStartup) Option {
	return func(o *option) {
		o.startupHandle = handle
	}
}

func WithShutdownHandle(handle OnShutdown) Option {
	return func(o *option) {
		o.shutdownHandle = handle
	}
}

func WithGormPlugins(plugins ...gorm.Plugin) Option {
	return func(o *option) {
		o.plugins = append(o.plugins, plugins...)
	}
}

func WithModels(models ...interface{}) Option {
	return func(o *option) {
		o.models = append(o.models, models...)
	}
}

func WithMgrAuth(auth ...app.HandlerFunc) Option {
	return func(o *option) {
		o.mgrAuth = append(o.mgrAuth, auth...)
	}
}

func WithValidators(validators map[string]func(args ...interface{}) error) Option {
	return func(o *option) {
		for k, v := range validators {
			o.validators[k] = v
		}
	}
}
