package main

import (
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart/utility"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestYamlConf(t *testing.T) {
	ym, _ := yaml.Marshal(Configuration{
		RDB:  utility.RdbProperties{},
		Cors: &Cors{},
	})
	t.Log(string(ym))
	//
	t.Log(sonic.MarshalString(model.PassportBinding{}))
}
