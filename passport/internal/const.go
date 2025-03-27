package internal

import (
	"github.com/ywengineer/smart-kit/pkg/app"
	"strconv"
)

var (
	// ErrTodo common error
	ErrTodo          = app.ApiError("common.err.todo_feature")
	ErrDisLock       = app.ApiError("common.err.lock")
	ErrCache         = app.ApiError("common.err.cache")
	ErrRdb           = app.ApiError("common.err.rdb")
	ErrJsonMarshal   = app.ApiError("common.err.json_marshal")
	ErrJsonUnmarshal = app.ApiError("common.err.json_unmarshal")
	ErrGenToken      = app.ApiError("common.err.gen_token")
	ErrInvalidToken  = app.ApiError("common.err.invalid_token")
	ErrBoundOther    = app.ApiError("common.err.bound_other")
	ErrPassword      = app.ApiError("common.err.passwd")
	ErrUserNotFound  = app.ApiError("common.err.user_not_found")
	ErrAuth          = app.ApiError("common.err.auth")
	ErrSign          = app.ApiError("common.err.sign")
	// ErrMaxPerDevice error for register service
	ErrMaxPerDevice = app.ApiError("register.err.max.account")
	ErrRegisterFail = app.ApiError("register.error.rdb")
	// ErrLoginTry error for login service
	ErrLoginTry = app.ApiError("login.err.asshole")
	// ErrSameBound
	ErrSameBound   = app.ApiError("bind.err.bound_type")
	ErrUnsupported = app.ApiError("bind.err.unsupported")
)

func ValidateErr(err error) app.ApiResult {
	if err == nil {
		return app.ApiError("validation.err", "ignore")
	}
	return app.ApiError("validation.err", err.Error())
}

func CacheKeyBoundTypes(passportId uint) string {
	return "bounds:" + strconv.FormatUint(uint64(passportId), 10)
}

func CacheKeyPassport(passportId uint) string {
	return "passport:" + strconv.FormatUint(uint64(passportId), 10)
}
