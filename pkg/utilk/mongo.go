package utilk

import (
	"context"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewMongo mongoUri format: see https://www.mongodb.com/docs/manual/reference/connection-string/
func NewMongo(ctx context.Context, mongoUri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(mongoUri).SetBSONOptions(&options.BSONOptions{
		NilByteSliceAsEmpty: true,
		NilMapAsEmpty:       true,
		NilSliceAsEmpty:     true,
		UseJSONStructTags:   true,
	})
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	go func() {
		_ = <-WatchContext(ctx)
		logk.DefaultLogger().Infof("mongo client disconnected: %s, error: %v", opts.Hosts, client.Disconnect(context.TODO()))
	}()
	return client, nil
}
