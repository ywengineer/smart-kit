package rpcs

import (
	"context"
	"testing"
	"time"
)

func TestRpc(t *testing.T) {
	n := make(chan any)
	GetDefaultRpc().GetAsync(context.Background(), "https://www.baidu.com", nil, func(statusCode int, body []byte, err error) {
		t.Logf(" code = %d\n body = %s\n err = %v", statusCode, string(body), err)
		n <- true
	})
	defer close(n)
	t.Logf("%v", <-n)
}

func TestTime(t *testing.T) {
	t.Log(time.Now().Format("2006-01-02T15:04:05.000Z07:00"))
	t.Log(time.Now().Format(time.RFC3339))
	t.Log(time.Now().UTC().Format(time.RFC3339))
}
