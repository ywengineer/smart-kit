package rpcs

import (
	"context"
	"io"
	"time"
)

type asyncRpc struct {
	t Rpc
}

func (h *asyncRpc) GetAsync(ctx context.Context, url string, callback RpcCallback) {
	rpcPool.CtxGo(ctx, func() {
		callback(h.t.Get(ctx, url))
	})
}

func (h *asyncRpc) GetTimeoutAsync(ctx context.Context, url string, timeout time.Duration, callback RpcCallback) {
	rpcPool.CtxGo(ctx, func() {
		callback(h.t.GetTimeout(ctx, url, timeout))
	})
}

func (h *asyncRpc) PostAsync(ctx context.Context, contentType string, url string, reqBody io.WriterTo, callback RpcCallback) {
	rpcPool.CtxGo(ctx, func() {
		callback(h.t.Post(ctx, contentType, url, reqBody))
	})
}
