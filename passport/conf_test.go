package main

import (
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
}
