package loaders

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type envMap struct {
	*Conf
	GoPath string `env:"GOPATH" json:"go_path"`
	User   string `env:"USER" json:"user"`
	Home   string `env:"HOME" json:"home"`
	Term   string `env:"TERM"`
}

func Test_EnvLoader(t *testing.T) {
	em := envMap{}
	err := NewLocalLoader("./conf.json").Load(&em)
	assert.Nil(t, err)
	err = NewEnvLoader().Load(&em)
	t.Logf("%#v\n", em)
	assert.Nil(t, err)
}
