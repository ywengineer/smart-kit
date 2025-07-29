package loaders

import (
	"context"
	"errors"
	"github.com/Netflix/go-env"
)

type envLoader struct {
}

func NewEnvLoader() SmartLoader {
	return &envLoader{}
}

func (ll *envLoader) Unmarshal(data []byte, out interface{}) error {
	_, err := env.UnmarshalFromEnviron(out)
	return err
}

func (ll *envLoader) Load(out interface{}) error {
	return ll.Unmarshal(nil, out)
}

func (ll *envLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return errors.New("env loader not support watch")
}
