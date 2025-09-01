package loaders

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

type ConfigType string

const (
	Yaml ConfigType = "yaml"
	Json            = "json"
)

type viperLoader struct {
	v  *viper.Viper
	ct ConfigType
}

func NewViperLoader(fileName string, configType ConfigType) SmartLoader {
	v := viper.New()
	v.SetConfigType(string(configType))
	v.SetConfigFile(fileName)
	v.AddConfigPath(".")
	v.AutomaticEnv()                                   // 启用环境变量替换（关键：将 ${ENV} 替换为实际环境变量）
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 允许使用 . 分隔符（如将 config 中的 db.host 对应环境变量 DB_HOST）
	return &viperLoader{v: v, ct: configType}
}

func (ll *viperLoader) Unmarshal(_ []byte, out interface{}) error {
	if err := ll.v.ReadInConfig(); err != nil {
		return err
	}
	return ll.v.Unmarshal(out)
}

func (ll *viperLoader) Load(out interface{}) error {
	return ll.Unmarshal(nil, out)
}

func (ll *viperLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return errors.New("viper loader not support watch")
}
