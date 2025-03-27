package signs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
)

// GenerateSignature 生成签名
func GenerateSignature(params map[string]string, secretKey string) string {
	cnt := len(params)
	if cnt == 0 {
		return ""
	}
	// 提取参数名并排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 按参数名自然顺序拼接参数值
	var paramStr string
	for ix, k := range keys {
		paramStr += url.QueryEscape(k) + "=" + url.QueryEscape(params[k])
		if ix+1 < cnt {
			paramStr += "&"
		}
	}
	// 创建一个新的 HMAC 哈希对象
	h := hmac.New(sha256.New, []byte(secretKey))
	// 写入待签名的数据
	h.Write([]byte(paramStr))
	// 计算哈希值
	mac := h.Sum(nil)
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(mac)
}

// VerifySignature 验证签名
func VerifySignature(params map[string]string, signature string, secretKey string) bool {
	// 去除签名参数
	if _, ok := params["signature"]; ok {
		delete(params, "signature")
	}
	// 比较签名
	return hmac.Equal([]byte(GenerateSignature(params, secretKey)), []byte(signature))
}
