package apps

import (
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type option struct {
	plugins []gorm.Plugin
	models  []interface{}
	mgrAuth []app.HandlerFunc
}

type Option func(*option)

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
