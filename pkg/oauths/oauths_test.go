package oauths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOauth(t *testing.T) {
	o := Oauth(map[string]map[string]string{
		"wx-id-1":    {"app-id": "app-id", "app-secret": "app-secret"},
		"wx-id-2":    {"app_id": "app-id", "app-secret": "app-secret"},
		"qq-id-1":    {"app-id": "app-id", "app-secret": "app-secret"},
		"qq-id-2":    {"app-id": "app-id", "app-secret": "app-secret"},
		"weibo-id-2": {"app-id": "app-id", "app-secret": "app-secret"},
	})
	//
	af, err := o.Get("wx-id-2")
	assert.NotNil(t, err)
	af, err = o.Get("wx-id-1")
	assert.Nil(t, err)
	_, err = af.GetToken("weibo-id-2")
	assert.NotNil(t, err)
	t.Logf("%v", err)
}
