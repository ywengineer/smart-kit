package internal

import (
	"github.com/ywengineer/smart-kit/pkg/apps"
	"strconv"
)

type UserType int

const (
	UserTypePlayer UserType = iota
	UserTypeMgr
)

const (
	TokenKeyUserType = "tkut"
)

var (
	// ErrTodo common error
	ErrTodo          = apps.ApiError("common.err.todo_feature")
	ErrDisLock       = apps.ApiError("common.err.lock")
	ErrCache         = apps.ApiError("common.err.cache")
	ErrRdb           = apps.ApiError("common.err.rdb")
	ErrJsonMarshal   = apps.ApiError("common.err.json_marshal")
	ErrJsonUnmarshal = apps.ApiError("common.err.json_unmarshal")
	ErrGenToken      = apps.ApiError("common.err.gen_token")
	ErrInvalidToken  = apps.ApiError("common.err.invalid_token")
	ErrBoundOther    = apps.ApiError("common.err.bound_other")
	ErrPassword      = apps.ApiError("common.err.passwd")
	ErrUserNotFound  = apps.ApiError("common.err.user_not_found")
	ErrAuth          = apps.ApiError("common.err.auth")
	ErrSign          = apps.ApiError("common.err.sign")
	// ErrMaxPerDevice error for register service
	ErrMaxPerDevice = apps.ApiError("register.err.max.account")
	ErrRegisterFail = apps.ApiError("register.error.rdb")
	// ErrLoginTry error for login service
	ErrLoginTry = apps.ApiError("login.err.asshole")
	// ErrSameBound
	ErrSameBound   = apps.ApiError("bind.err.bound_type")
	ErrUnsupported = apps.ApiError("bind.err.unsupported")
)

func ValidateErr(err error) apps.ApiResult {
	if err == nil {
		return apps.ApiError("validation.err", "ignore")
	}
	return apps.ApiError("validation.err", err.Error())
}

func CacheKeyBoundTypes(passportId uint) string {
	return "bounds:" + strconv.FormatUint(uint64(passportId), 10)
}

func CacheKeyPassport(passportId uint) string {
	return "passport:" + strconv.FormatUint(uint64(passportId), 10)
}
