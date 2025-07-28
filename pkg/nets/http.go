package nets

import "github.com/cloudwego/hertz/pkg/protocol/consts"

// Is2xx 用于检查状态码是否为 2xx
func Is2xx(statusCode int) bool {
	return statusCode >= consts.StatusOK && statusCode < consts.StatusMultipleChoices
}
