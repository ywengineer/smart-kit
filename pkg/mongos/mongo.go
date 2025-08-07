package mongos

import (
	"context"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var define = make(map[Collection][]CollectionIndex)

type Collection string

type CollectionIndex struct {
	Name               string // 索引名称（可选，不填则由MongoDB自动生成）
	Keys               bson.D // 索引键，如：{"email": 1, "age": -1}（1升序，-1降序）
	Unique             bool   // 是否唯一索引
	ExpireAfterSeconds *int32 // 其他可选参数（如过期时间等）仅适用于TTL索引
}

func Register(collection Collection, index CollectionIndex) {
	define[collection] = append(define[collection], index)
}

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
		_ = <-utilk.WatchContext(ctx)
		logk.Infof("mongo client disconnected: %s, error: %v", opts.Hosts, client.Disconnect(context.TODO()))
	}()
	return client, nil
}

func EnsureCollection(ctx context.Context, db *mongo.Database, extra map[Collection][]CollectionIndex) error {
	if len(extra) > 0 {
		for collection, indexes := range extra {
			define[collection] = append(define[collection], indexes...)
		}
	}
	c, err := db.ListCollections(ctx, bson.D{})
	if err != nil {
		return err
	}
	defer c.Close(ctx)
	// ----------------------------------------------------------------------------------------------------
	var collList []string
	for c.Next(ctx) {
		name := c.Current.Lookup("name").StringValue()
		collList = append(collList, name)
	}
	if err = c.Err(); err != nil {
		return err
	}
	// ----------------------------------------------------------------------------------------------------
	for coll, indexes := range define {
		strColl := string(coll)
		if !lo.Contains(collList, strColl) {
			if err = db.CreateCollection(ctx, strColl); err != nil {
				return err
			}
		}
		if err = ensureIndexes(ctx, db.Collection(strColl), indexes); err != nil {
			return err
		}
	}
	return nil
}

// ensureIndexes 初始化集合的索引（若不存在则创建）
func ensureIndexes(ctx context.Context, coll *mongo.Collection, indexes []CollectionIndex) error {
	// existing indexes
	cursor, err := coll.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	// -----------------------------------------------------------------------------------------------
	var existsIndexes []string
	for cursor.Next(ctx) {
		existsIndexes = append(existsIndexes, cursor.Current.Lookup("name").String())
	}
	if err = cursor.Err(); err != nil {
		return err
	}
	// -----------------------------------------------------------------------------------------------
	for _, idx := range indexes {
		// skip exists index
		if lo.Contains(existsIndexes, idx.Name) {
			continue
		}
		// create index options
		opts := options.Index().SetUnique(idx.Unique)
		if len(idx.Name) > 0 {
			opts = opts.SetName(idx.Name)
		}
		if idx.ExpireAfterSeconds != nil {
			opts = opts.SetExpireAfterSeconds(*idx.ExpireAfterSeconds)
		}
		// create index
		if _, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    idx.Keys,
			Options: opts,
		}); err != nil {
			return err
		}
	}
	// -----------------------------------------------------------------------------------------------
	return nil
}
