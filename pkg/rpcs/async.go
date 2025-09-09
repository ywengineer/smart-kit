package rpcs

import (
	"context"
	"io"
	"net/http"
)

type asyncRpc struct {
	t Rpc
}

func (h *asyncRpc) GetAsync(ctx context.Context, url string, header http.Header, callback RpcCallback) {
	rpcPool.CtxGo(ctx, func() {
		callback(h.t.Get(ctx, url, header))
	})
}

func (h *asyncRpc) PostAsync(ctx context.Context, contentType string, url string, header http.Header, reqBody io.WriterTo, callback RpcCallback) {
	rpcPool.CtxGo(ctx, func() {
		callback(h.t.Post(ctx, contentType, url, header, reqBody))
	})
}
