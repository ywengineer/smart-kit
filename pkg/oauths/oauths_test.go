package oauths

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWxOauth(t *testing.T) {
	o := Oauth(map[string]map[string]string{
		"wx-id-1": {"app-id": "app-id", "app-secret": "app-secret", "type": "wx"},
		"wx-id-2": {"app_id": "app-id", "app-secret": "app-secret", "type": "wx"},
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

func TestQQOauth(t *testing.T) {
	o := Oauth(map[string]map[string]string{
		"qq-id-1": {"app-id": "app-id", "app-secret": "app-secret", "type": "qq", "redirect-url": "redirect-url"},
		"qq-id-2": {"app_did": "app-id", "app-secret": "app-secret", "type": "qq", "redirect-url": "redirect-url"},
	})
	//
	af, err := o.Get("qq-id-2")
	assert.NotNil(t, err)
	t.Logf("%v", err)
	af, err = o.Get("qq-id-1")
	assert.Nil(t, err)
	_, err = af.GetToken("auth_code")
	assert.NotNil(t, err)
	t.Logf("%v", err)
	_, err = af.GetUserInfo("", "access_token")
	assert.NotNil(t, err)
	t.Logf("%v", err)
}
