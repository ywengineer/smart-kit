package loaders

import (
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	ServiceName       string                 `json:"service_name" yaml:"service-name" mapstructure:"service-name" env:"SERVICE_NAME"`
	Weight            int                    `json:"weight" yaml:"weight" mapstructure:"weight" env:"SERVICE_WEIGHT"`
	Network           string                 `json:"network" yaml:"network" mapstructure:"network" env:"SERVICE_NETWORK"`
	Address           string                 `json:"address" yaml:"address" mapstructure:"address" env:"SERVICE_ADDRESS"`
	Workers           int                    `json:"workers" yaml:"workers" mapstructure:"workers" env:"SERVICE_WORKERS"`
	WorkerLoadBalance string                 `json:"load_balance" yaml:"load-balance" mapstructure:"load-balance" env:"SERVICE_LOAD_BALANCE"`
	Metadata          map[string]interface{} `json:"metadata" yaml:"metadata" mapstructure:"metadata" env:"SERVICE_METADATA"`
}

func (c *Conf) String() string {
	s, _ := sonic.MarshalString(c)
	return s
}

func (c *Conf) Get(key string, def string) string {
	if v, ok := c.Metadata[key]; ok {
		return utilk.ToString(v)
	}
	return def
}

type Decoder interface {
	Unmarshal(buf []byte, val interface{}) error
}

func NewJsonDecoder() Decoder {
	return &jsonDecoder{}
}

func NewYamlDecoder() Decoder {
	return &yamlDecoder{}
}

type jsonDecoder struct {
}

func (d *jsonDecoder) Unmarshal(buf []byte, val interface{}) error {
	return sonic.Unmarshal(buf, val)
}

type yamlDecoder struct {
}

func (d *yamlDecoder) Unmarshal(buf []byte, val interface{}) error {
	return yaml.Unmarshal(buf, val)
}
