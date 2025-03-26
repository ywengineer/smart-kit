package rpcs

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/discovery"
	"github.com/cloudwego/hertz/pkg/app/client/retry"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	client_http "github.com/cloudwego/hertz/pkg/protocol/client"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/ywengineer/smart-kit/pkg/utilk"
	"io"
	"time"
)

const defaultMaxRedirectsCount = 16

func NewDefaultRpc(info RpcClientInfo) (rpc Rpc, err error) {
	return NewHertzRpc(nil, info)
}

func NewHertzRpc(resolver discovery.Resolver, info RpcClientInfo) (rpc Rpc, err error) {
	var cli *client.Client
	info.MaxRetry = utilk.MaxInt(1, info.MaxRetry)
	info.MaxConnPerHost = utilk.MaxInt(50, info.MaxConnPerHost)
	if info.Delay < time.Millisecond*50 {
		info.Delay = time.Millisecond * 50
	}
	if cli, err = client.NewClient(
		client.WithMaxConnsPerHost(info.MaxConnPerHost),
		client.WithName(info.ClientName),
		client.WithClientReadTimeout(time.Second*5),
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
			hlog.Debugf("[RPC][%s] [cost %dms] invoke target: %s ", info.ClientName, time.Now().Unix()-ts, req.RequestURI())
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

func (h *hertzRPC) Get(ctx context.Context, url string) (statusCode int, body []byte, err error) {
	return h.cli.Get(ctx, nil, url, h.cluster)
}

func (h *hertzRPC) GetTimeout(ctx context.Context, url string, timeout time.Duration) (statusCode int, body []byte, err error) {
	return h.cli.GetTimeout(ctx, nil, url, timeout, h.cluster)
}

// Post contentType see consts.MIMEXXX
func (h *hertzRPC) Post(ctx context.Context, contentType string, url string, reqBody io.WriterTo) (statusCode int, body []byte, err error) {
	req := protocol.AcquireRequest()
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte(contentType))
	req.SetOptions(h.cluster)
	//
	if reqBody != nil {
		if _, err := reqBody.WriteTo(req.BodyWriter()); err != nil {
			return 0, nil, err
		}
	}
	//
	statusCode, body, err = h.doRequestFollowRedirectsBuffer(ctx, req, nil, url)
	//
	protocol.ReleaseRequest(req)
	return statusCode, body, err
}

func (h *hertzRPC) doRequestFollowRedirectsBuffer(ctx context.Context, req *protocol.Request, dst []byte, url string) (statusCode int, body []byte, err error) {
	resp := protocol.AcquireResponse()
	bodyBuf := resp.BodyBuffer()
	oldBody := bodyBuf.B
	bodyBuf.B = dst

	statusCode, _, err = client_http.DoRequestFollowRedirects(ctx, req, resp, url, defaultMaxRedirectsCount, h.cli)

	// In HTTP2 scenario, client use stream mode to create a request and its body is in body stream.
	// In HTTP1, only client recv body exceed max body size and client is in stream mode can trig it.
	body = resp.Body()
	bodyBuf.B = oldBody
	protocol.ReleaseResponse(resp)

	return statusCode, body, err
}
