package apps

import "gorm.io/gorm"

type option struct {
	plugins []gorm.Plugin
	models  []interface{}
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
