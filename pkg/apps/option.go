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
	middlewares    []app.HandlerFunc
}

type Option func(*option)

func WithMiddlewares(middlewares ...app.HandlerFunc) Option {
	return func(o *option) {
		o.middlewares = append(o.middlewares, middlewares...)
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
