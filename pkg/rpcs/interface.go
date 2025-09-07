package rpcs

import (
	"context"
	"io"
	"time"

	"github.com/bytedance/gopkg/util/gopool"
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/proto"
)

const (
	ContentTypeUrlencoded = "application/x-www-form-urlencoded"
	ContentTypeFormData   = "multipart/form-data"
	ContentTypeJSON       = "application/json"
	ContentTypeProtoBuf   = "application/x-protobuf"
)

var rpcPool = gopool.NewPool("rpc-pool", 500, gopool.NewConfig())

type RpcClientInfo struct {
	ClientName     string        `json:"client_name" yaml:"client-name"`
	MaxRetry       uint          `json:"max_retry" yaml:"max-retry"`
	Delay          time.Duration `json:"retry-delay" yaml:"retry-delay"`
	MaxConnPerHost int           `json:"max_conn_per_host" yaml:"max-conn-per-host"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read-timeout"`
}

type RpcCallback func(statusCode int, body []byte, err error)

type Rpc interface {
	Get(ctx context.Context, url string) (statusCode int, body []byte, err error)

	GetTimeout(ctx context.Context, url string, timeout time.Duration) (statusCode int, body []byte, err error)

	Post(ctx context.Context, contentType string, url string, reqBody io.WriterTo) (statusCode int, body []byte, err error)

	GetAsync(ctx context.Context, url string, callback RpcCallback)

	GetTimeoutAsync(ctx context.Context, url string, timeout time.Duration, callback RpcCallback)

	PostAsync(ctx context.Context, contentType string, url string, reqBody io.WriterTo, callback RpcCallback)
}

type JsonBody struct {
	V any
}

func (j JsonBody) WriteTo(w io.Writer) (n int64, err error) {
	var b []byte
	var wn = 0
	if b, err = sonic.Marshal(j.V); err == nil {
		wn, err = w.Write(b)
	}
	return int64(wn), err
}

type ProtoBody struct {
	V proto.Message
}

func (p ProtoBody) WriteTo(w io.Writer) (n int64, err error) {
	var b []byte
	var wn = 0
	if b, err = proto.Marshal(p.V); err == nil {
		wn, err = w.Write(b)
	}
	return int64(wn), err
}
