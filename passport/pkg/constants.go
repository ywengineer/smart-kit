package pkg

const (
	ContextKeySmart   = "smart-context"
	ContextKeyRDB     = "rdb"
	ContextKeyRedis   = "redis"
	ContextKeyApp     = "app"
	ContextKeyRedLock = "red-lock"
)

type ApiResult struct {
	Code    ResultCode  `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ResultCode int

const (
	Fail                ResultCode = 0
	Ok                             = 200
	InternalServerError            = 503
)

func ApiError(msg string, data ...interface{}) ApiResult {
	return ApiResult{
		Code:    InternalServerError,
		Data:    _first(data),
		Message: msg,
	}
}

func ApiFail(msg string, data ...interface{}) ApiResult {
	return ApiResult{
		Code:    Fail,
		Data:    _first(data),
		Message: msg,
	}
}

func ApiOk(data interface{}) ApiResult {
	return ApiResult{
		Code:    Ok,
		Data:    data,
		Message: "ok",
	}
}

func _first(arr []interface{}) interface{} {
	if arr == nil || len(arr) == 0 {
		return nil
	}
	return arr[0]
}
