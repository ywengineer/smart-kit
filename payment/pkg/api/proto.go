package api

import (
	"gitee.com/ywengineer/smart-kit/payment/pkg/proto"
)

type ProtoErrCode int32

func (p ProtoErrCode) Code() int32 {
	return int32(p)
}

const (
	C0     ProtoErrCode = 0     // OK.
	C21000              = 21000 // The App Store could not read the JSON object you provided.
	C21002              = 21002 // The data in the receipt-data property was malformed or missing.
	C21003              = 21003 // The receipt could not be authenticated.
	C21004              = 21004 // The shared secret you provided does not match the shared secret on file for your account.
	C21005              = 21005 // The receipt server is not currently available.
	C21006              = 21006 // This receipt is valid but the subscription has expired. When this status code is returned to your server, the receipt data is also decoded and returned as part of the response. Only returned for iOS 6 style transaction receipts for auto-renewable
	// subscriptions.
	C21007       = 21007 // This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead.
	C21008       = 21008 // This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead.
	C21010       = 21010 // This receipt could not be authorized. Treat this the same as if a purchase was never made.
	C21100_21199 = 21100 // 21100 - 21199 Internal data access error.
	// app
	C90000 = 90000 // 订单已处理
	C90001 = 90001 // 无效product_id
	C90002 = 90002 // 无效bundle_id
	C90003 = 90003 // 订阅过期
	C90004 = 90004 // 订阅续订账号不匹配
)

func NewProtoOkResult(message string) *proto.ApiResult {
	return NewProtoResult(c01, message, C0)
}

func NewProtoExceptionResult(err error, errCode ProtoErrCode) *proto.ApiResult {
	return NewProtoResult(c02, err.Error(), errCode)
}

func NewProtoFailCodeResult(errCode ProtoErrCode) *proto.ApiResult {
	return NewProtoResult(c03, "", errCode)
}

func NewProtoFailResult(message string, errCode ProtoErrCode) *proto.ApiResult {
	return NewProtoResult(c03, message, errCode)
}

func NewProtoResult(code Code, message string, errCode ProtoErrCode) *proto.ApiResult {
	r := new(proto.ApiResult)
	c, ec := string(code), errCode.Code()
	r.Code = &c
	r.Msg = &message
	r.ErrorCode = &ec
	return r
}
