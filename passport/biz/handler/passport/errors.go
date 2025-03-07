package passport

import "github.com/ywengineer/smart-kit/passport/pkg"

var (
	// ErrTodo common error
	ErrTodo          = pkg.ApiError("common.err.todo_feature")
	ErrDisLock       = pkg.ApiError("common.err.lock")
	ErrCache         = pkg.ApiError("common.err.cache")
	ErrRdb           = pkg.ApiError("common.err.rdb")
	ErrJsonMarshal   = pkg.ApiError("common.err.json_marshal")
	ErrJsonUnmarshal = pkg.ApiError("common.err.json_unmarshal")
	ErrGenToken      = pkg.ApiError("common.err.gen_token")
	ErrInvalidToken  = pkg.ApiError("common.err.invalid_token")
	// ErrMaxPerDevice error for register service
	ErrMaxPerDevice = pkg.ApiError("register.err.max.account")
	ErrBoundOther   = pkg.ApiError("register.bound_other")
	ErrRegisterFail = pkg.ApiError("register.error.rdb")
	// ErrLoginTry error for login service
	ErrLoginTry = pkg.ApiError("login.err.asshole")
)

func validateErr(err error) pkg.ApiResult {
	if err == nil {
		return pkg.ApiError("validation.err", "ignore")
	}
	return pkg.ApiError("validation.err", err.Error())
}
