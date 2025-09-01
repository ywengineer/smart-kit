package loaders

import (
	"github.com/bytedance/sonic"
	"testing"
)

func TestViperYamlLoader(t *testing.T) {
	//
	loader := NewViperLoader("conf", Yaml)
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Log(sonic.MarshalString(c))
	}
}
