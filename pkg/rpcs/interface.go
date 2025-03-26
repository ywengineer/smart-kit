package rpcs

import (
	"context"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/bytedance/sonic"
	"io"
	"math"
	"time"
)

const (
	ContentTypeUrlencoded = "application/x-www-form-urlencoded"
	ContentTypeFormData   = "multipart/form-data"
	ContentTypeJSON       = "application/json"
)

var rpcPool = gopool.NewPool("rpc-pool", math.MaxInt, gopool.NewConfig())

type RpcClientInfo struct {
	ClientName     string        `json:"client_name" yaml:"client-name"`
	MaxRetry       uint          `json:"max_retry" yaml:"max-retry"`
	Delay          time.Duration `json:"retry-delay" yaml:"retry-delay"`
	MaxConnPerHost int           `json:"max_conn_per_host" yaml:"max-conn-per-host"`
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

func NewJSONBody(v any) io.WriterTo {
	return &jsonBody{v: v}
}

type jsonBody struct {
	v any
}

func (j *jsonBody) WriteTo(w io.Writer) (n int64, err error) {
	var b []byte
	var wn = 0
	if b, err = sonic.Marshal(j.v); err == nil {
		wn, err = w.Write(b)
	}
	return int64(wn), err
}
