package rpcs

import (
	"context"
	"testing"
)

func TestRpc(t *testing.T) {
	n := make(chan any)
	GetDefaultRpc().GetAsync(context.Background(), "https://www.baidu.com", func(statusCode int, body []byte, err error) {
		t.Logf(" code = %d\n body = %s\n err = %v", statusCode, string(body), err)
		n <- true
	})
	defer close(n)
	t.Logf("%v", <-n)
}
