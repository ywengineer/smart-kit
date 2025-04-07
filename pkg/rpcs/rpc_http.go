package rpcs

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"io"
	"net/http"
	"sync"
)

var standard = &http.Client{
	Transport: &RoundTripper{rpc: defaultRpc},
}

func StandardBasedOnHertz() *http.Client {
	return standard
}

// RoundTripper implements the http.RoundTripper interface, using a rpc client to execute requests.
type RoundTripper struct {
	// The client to use during requests. If nil, the default Rpc and settings will be used.
	rpc Rpc
	// once ensures that the logic to initialize the default client runs at
	// most once, in a single thread.
	once sync.Once
}

// init initializes the underlying rpc client.
func (rt *RoundTripper) init() {
	if rt.rpc == nil {
		rt.rpc = defaultRpc
	}
}

// RoundTrip satisfies the http.RoundTripper interface.
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.once.Do(rt.init)
	//
	reqh := protocol.AcquireRequest()
	reqh.SetOptions(config.WithSD(false))
	defer protocol.ReleaseRequest(reqh)
	//
	if err := convertHTTPToHertz(req, reqh); err != nil {
		return nil, err
	}
	var headers = make(http.Header)
	//
	if statusCode, body, err := rt.rpc.(*hertzRPC).doRequestFollowRedirectsBuffer2(context.Background(), reqh, nil, req.URL.RequestURI(), headers); err != nil {
		return nil, err
	} else {
		httpResp := &http.Response{
			Status:        http.StatusText(statusCode),
			StatusCode:    statusCode,
			Header:        headers,
			ContentLength: int64(len(body)),
			Body:          http.NoBody,
		}
		if httpResp.ContentLength > 0 {
			httpResp.Body = io.NopCloser(bytes.NewReader(body))
		}
		return httpResp, nil
	}
}

func convertHTTPToHertz(httpReq *http.Request, hertzReq *protocol.Request) error {
	// 复制基础信息
	hertzReq.Header.SetMethod(httpReq.Method)
	hertzReq.SetRequestURI(httpReq.URL.RequestURI())
	hertzReq.ParseURI()
	hertzReq.SetHost(httpReq.Host)
	// 复制 Headers
	for key, values := range httpReq.Header {
		for _, value := range values {
			hertzReq.Header.Add(key, value)
		}
	}
	// 复制 Query 参数
	for key, values := range httpReq.URL.Query() {
		for _, value := range values {
			hertzReq.URI().QueryArgs().Add(key, value)
		}
	}
	// 复制 Form 数据（需先 ParseForm）
	if err := httpReq.ParseForm(); err == nil {
		for key, values := range httpReq.Form {
			for _, value := range values {
				hertzReq.PostArgs().Add(key, value)
			}
		}
	}
	// 处理 Body
	if httpReq.Body != nil {
		body, err := io.ReadAll(httpReq.Body)
		if err != nil {
			return fmt.Errorf("read body failed: %w", err)
		}
		hertzReq.SetBody(body)
		httpReq.Body = io.NopCloser(bytes.NewReader(body))
	}
	return nil
}
