package hw

import (
	"os"
	"testing"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	//
	keyId := os.Getenv("HW_APP_ID")
	pKey := os.Getenv("HW_APP_SECRET")
	assert.NotEmpty(t, keyId)
	assert.NotEmpty(t, pKey)
	t.Log(keyId)
	t.Log(pKey)
	cf := config.ChannelProperty{
		Validator:    "huawei",
		ClientID:     keyId, // 从控制台获取
		ClientSecret: pKey,  // 从控制台获取（严格保密）
		Sandbox:      true,  // 测试用 true，生产用 false
		ApiRoot:      "https://orders-drru.iap.cloud.huawei.ru",
		Apps:         []string{},
	}
	r, err := New(cf)
	assert.Nil(t, err)
	//
	tk, err := r.Verify(t.Context(), "123456")
	assert.Nil(t, err)
	t.Logf("Result: %+v", tk)
}
