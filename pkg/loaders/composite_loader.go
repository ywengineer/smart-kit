package loaders

import (
	"context"
)

type compLoader struct {
	loaderContainer []SmartLoader
}

func NewCompositeLoader(loaders ...SmartLoader) SmartLoader {
	return &compLoader{loaderContainer: loaders}
}

func (ll *compLoader) Unmarshal(data []byte, out interface{}) error {
	return nil
}

func (ll *compLoader) Load(out interface{}) error {
	for _, loader := range ll.loaderContainer {
		if err := loader.Load(out); err != nil {
			return err
		}
	}
	return nil
}

// Watch only support watch the last SmartLoader
func (ll *compLoader) Watch(ctx context.Context, callback WatchCallback) error {
	return ll.loaderContainer[len(ll.loaderContainer)-1].Watch(ctx, callback)
}
