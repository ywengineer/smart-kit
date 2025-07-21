package utilk

import (
	"context"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"os"
	"testing"
	"time"
)

type Parent struct {
	Id       int64     `json:"_id"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	DeleteAt time.Time `json:"delete_at"`
}

type OrmObject struct {
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	Int       int64      `json:"int"`
	OrmObject *OrmObject `json:"orm_object,omitempty"`
}

var ctx = context.Background()
var cancel context.CancelFunc
var coll *mongo.Collection

func setup() {
	ctx, cancel = context.WithCancel(ctx)
	db, err := NewMongo(ctx, "mongodb://localhost:27017/test")
	if err != nil {
		panic(err)
	}
	coll = db.Database("test").Collection("test")
}

func teardown() {
	cancel()
	time.Sleep(3 * time.Second)
}

func TestMain(m *testing.M) {
	// 初始化（如连接数据库、启动服务器）
	setup()
	// 执行所有测试
	code := m.Run()
	// 清理资源
	teardown()
	os.Exit(code)
}

func TestNewMongo(t *testing.T) {
	_, err := coll.InsertOne(ctx, OrmObject{Name: "test:" + rand.Text(), Password: "testp", Int: 1})
	assert.Nil(t, err)
	_, err = coll.InsertOne(ctx, OrmObject{Name: "test:" + rand.Text(), Password: "testp", Int: 1, OrmObject: &OrmObject{Name: "test:" + rand.Text(), Password: "testpp"}})
	assert.Nil(t, err)
}

func TestFindOne(t *testing.T) {
	oo := OrmObject{Name: "test"}
	err := coll.FindOne(ctx, bson.M{"name": "test"}).Decode(&oo)
	assert.Nil(t, err)
	t.Logf("TestFindOne: %+v", oo)
}
