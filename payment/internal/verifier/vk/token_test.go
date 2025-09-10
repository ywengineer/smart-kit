package vk

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTokenManager(t *testing.T) {
	//
	keyId := os.Getenv("VK_APP_ID")
	pKey := os.Getenv("VK_APP_SECRET")
	assert.NotEmpty(t, keyId)
	assert.NotEmpty(t, pKey)
	t.Log(keyId)
	t.Log(pKey)
	config := RustoreConfig{
		ClientID:     keyId, // 从控制台获取
		ClientSecret: pKey,  // 从控制台获取（严格保密）
		IsSandbox:    true,  // 测试用 true，生产用 false
		Apps:         []string{},
	}
	tm, err := NewTokenManager(config)
	assert.Nil(t, err)
	//
	tk, err := tm.getToken()
	assert.Nil(t, err)
	t.Logf("Token: %v", tk)
}
