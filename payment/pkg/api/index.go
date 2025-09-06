package api

import (
	"gitee.com/ywengineer/smart-kit/payment/biz/model/common"
)

type Code string

const (
	c01 = "c01" // 正常处理成功
	c02 = "c02" // 服务器异常
	c03 = "c03" // 处理失败(错误消息)
)

type ErrCode string

const None ErrCode = ""

func NewOkResult(message string) *common.ApiResult {
	return NewResult(c01, message, None)
}

func NewExceptionResult(err error, errCode ErrCode) *common.ApiResult {
	return NewResult(c02, err.Error(), errCode)
}

func NewFailResult(message string, errCode ErrCode) *common.ApiResult {
	return NewResult(c03, message, errCode)
}

func NewResult(code Code, message string, errCode ErrCode) *common.ApiResult {
	r := common.NewApiResult()
	r.Code = string(code)
	r.Message = message
	r.ErrCode = string(errCode)
	return r
}
