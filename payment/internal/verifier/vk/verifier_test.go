package vk

import (
	"os"
	"testing"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	//
	keyId := os.Getenv("VK_APP_ID")
	pKey := os.Getenv("VK_APP_SECRET")
	assert.NotEmpty(t, keyId)
	assert.NotEmpty(t, pKey)
	t.Log(keyId)
	t.Log(pKey)
	cf := config.ChannelProperty{
		Validator:    "rustore",
		ClientID:     keyId, // 从控制台获取
		ClientSecret: pKey,  // 从控制台获取（严格保密）
		Sandbox:      true,  // 测试用 true，生产用 false
		Apps:         []string{"2063556499"},
	}
	r, err := New(cf)
	assert.Nil(t, err)
	//
	tk, err := r.Verify(t.Context(), "10000015096")
	assert.Nil(t, err)
	t.Logf("Result: %+v", tk)
}
