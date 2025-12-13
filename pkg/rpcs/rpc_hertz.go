package rpcs

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"runtime"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/discovery"
	"github.com/cloudwego/hertz/pkg/app/client/retry"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	client_http "github.com/cloudwego/hertz/pkg/protocol/client"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const defaultMaxRedirectsCount = 16

var defaultRpc, _ = newDefaultRpc(RpcClientInfo{ClientName: "default-smart-rpc-client", ReadTimeout: 5 * time.Second})

// GetDefaultRpc retry = 1, delay = 50ms, read_timout = 100ms, max_retry = 1, max_conn_per_host = 100
func GetDefaultRpc() Rpc {
	return defaultRpc
}

func newDefaultRpc(info RpcClientInfo) (rpc Rpc, err error) {
	return NewHertzRpc(nil, info)
}

func NewHertzRpc(resolver discovery.Resolver, info RpcClientInfo) (rpc Rpc, err error) {
	var cli *client.Client
	info.MaxRetry = utilk.Max(1, info.MaxRetry)
	info.MaxConnPerHost = utilk.Max(runtime.GOMAXPROCS(0)+1, info.MaxConnPerHost)
	info.Delay = utilk.Max(info.Delay, time.Millisecond*1000)
	info.ReadTimeout = utilk.Max(info.ReadTimeout, time.Millisecond*100)
	if cli, err = client.NewClient(
		client.WithName(info.ClientName),
		client.WithMaxConnsPerHost(info.MaxConnPerHost),
		client.WithDialTimeout(5*time.Second),
		client.WithClientReadTimeout(info.ReadTimeout),
		client.WithMaxConnWaitTimeout(time.Second),
		client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		client.WithDialer(standard.NewDialer()),
		client.WithRetryConfig(retry.WithMaxAttemptTimes(info.MaxRetry), retry.WithDelayPolicy(retry.BackOffDelayPolicy)),
	); err != nil {
		return nil, err
	}
	if resolver != nil {
		cli.Use(sd.Discovery(resolver))
	}
	cli.Use(func(next client.Endpoint) client.Endpoint {
		return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
			ts := time.Now().Unix()
			n := next(ctx, req, resp)
			hlog.Debugf("[RPC][%s] [cost %dms] invoke [host=%s] [target=%s] [status=%d]", info.ClientName, time.Now().Unix()-ts, req.Host(), req.RequestURI(), resp.StatusCode())
			return n
		}
	})
	//
	return (&hertzRPC{cli: cli, cluster: config.WithSD(resolver != nil)}).init(), nil
}

type hertzRPC struct {
	*asyncRpc
	cli     *client.Client
	cluster config.RequestOption
}

func (h *hertzRPC) init() Rpc {
	h.asyncRpc = &asyncRpc{t: h}
	return h
}

func (h *hertzRPC) Get(ctx context.Context, url string, header http.Header) (statusCode int, body []byte, err error) {
	req := protocol.AcquireRequest()
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}
	req.SetOptions(h.cluster)

	statusCode, body, err = h.doRequestFollowRedirectsBuffer(ctx, req, nil, url)

	protocol.ReleaseRequest(req)
	return statusCode, body, err
}

// Post contentType see consts.MIMEXXX
func (h *hertzRPC) Post(ctx context.Context, contentType string, url string, header http.Header, reqBody io.WriterTo) (statusCode int, body []byte, err error) {
	req := protocol.AcquireRequest()
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte(contentType))
	req.SetOptions(h.cluster)
	//
	if reqBody != nil {
		if cl, err := reqBody.WriteTo(req.BodyWriter()); err != nil {
			return 0, nil, err
		} else {
			req.Header.SetContentLength(int(cl))
		}
	}
	//
	statusCode, body, err = h.doRequestFollowRedirectsBuffer(ctx, req, nil, url)
	//
	protocol.ReleaseRequest(req)
	return statusCode, body, err
}

func (h *hertzRPC) doRequestFollowRedirectsBuffer(ctx context.Context, req *protocol.Request, dst []byte, url string) (statusCode int, body []byte, err error) {
	statusCode, body, err = h.doRequestFollowRedirectsBuffer2(ctx, req, dst, url, nil)
	return statusCode, body, err
}

func (h *hertzRPC) doRequestFollowRedirectsBuffer2(ctx context.Context, req *protocol.Request, dst []byte, url string, respHeaders http.Header) (statusCode int, body []byte, err error) {
	resp := protocol.AcquireResponse()
	bodyBuf := resp.BodyBuffer()
	oldBody := bodyBuf.B
	bodyBuf.B = dst

	statusCode, _, err = client_http.DoRequestFollowRedirects(ctx, req, resp, url, defaultMaxRedirectsCount, h.cli)

	// In HTTP2 scenario, client use stream mode to create a request and its body is in body stream.
	// In HTTP1, only client recv body exceed max body size and client is in stream mode can trig it.
	body = resp.Body()
	bodyBuf.B = oldBody
	//
	if respHeaders != nil {
		resp.Header.VisitAll(func(k, v []byte) {
			respHeaders[string(k)] = []string{string(v)}
		})
	}
	resp.Header.GetProtocol()
	//
	protocol.ReleaseResponse(resp)
	return statusCode, body, err
}
