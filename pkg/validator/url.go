package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// 用于验证相对路径的正则表达式
// 允许字母、数字、斜杠、短横线、下划线、点
var relativePathRegex = regexp.MustCompile(`^/[a-zA-Z0-9/_\-.?&=]*$`)

var allowProtocols = []string{"http", "https"}

var ErrURLInvalid = errors.New("url is invalid")

// URLValidator 自定义 URL 校验器：支持协议白名单（默认 http/https）
func URLValidator(args ...interface{}) error {
	if len(args) <= 0 {
		return ErrURLInvalid
	}
	rawURL := args[0].(string)
	if rawURL == "" {
		return ErrURLInvalid // 空字符串由 required 标签控制，此处可忽略
	}
	// 1. 语法校验
	parsedURL, err := url.ParseRequestURI(rawURL)
	// 校验 URL 语法是否正确
	if err == nil && parsedURL.Scheme != "" {
		// 2. 协议校验（支持通过 tag 参数指定白名单，如 binding:"url(http,https,ftp)"）
		// 解析 tag 中的协议白名单（默认 http/https）
		protocol := strings.ToLower(parsedURL.Scheme)
		found := false
		for _, p := range allowProtocols {
			if protocol == p {
				found = true
				break
			}
		}
		if !found {
			return ErrURLInvalid
		}
		return nil
	}
	// 不是绝对URL，尝试按站内相对路径验证
	// 我们只允许以 '/' 开头的绝对路径，更安全可控
	if rawURL[0] == '/' {
		// 使用正则表达式验证路径格式
		if !relativePathRegex.MatchString(rawURL) {
			return ErrURLInvalid
		} else {
			return nil
		}
	}
	return ErrURLInvalid
}
