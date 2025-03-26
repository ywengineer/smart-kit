package oauths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOauth1(t *testing.T) {
	o := Oauth(map[string]map[string]string{
		"wx-id-1":    {"app-id": "app-id", "app-secret": "app-secret", "type": "wx"},
		"wx-id-2":    {"app_id": "app-id", "app-secret": "app-secret", "type": "wx"},
		"qq-id-1":    {"app-id": "app-id", "app-secret": "app-secret", "type": "qq"},
		"qq-id-2":    {"app-id": "app-id", "app-secret": "app-secret", "type": "qq"},
		"weibo-id-2": {"app-id": "app-id", "app-secret": "app-secret", "type": "weibo"},
	})
	//
	af, err := o.Get("wx-id-2")
	assert.NotNil(t, err)
	t.Logf("%v", err)
	af, err = o.Get("wx-id-1")
	assert.Nil(t, err)
	_, err = af.GetToken("weibo-id-2")
	assert.NotNil(t, err)
	t.Logf("%v", err)
	af, err = o.Get("weibo-id-2")
	assert.NotNil(t, err)
	t.Logf("%v", err)
}
