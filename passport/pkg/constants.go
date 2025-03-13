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

func (ar *ApiResult) Ok() bool {
	return ar.Code == Ok
}

type ResultCode int

const (
	DeviceId    string = "device_id"    // DeviceId
	DeviceModel string = "device_model" // DeviceModel 机型 (deviceModel)
	GameVersion string = "ver"          // GameVersion 游戏版本 (v)
	Os          string = "os"           // Os
	OsVersion   string = "os_ver"       // OsVersion 系统版本	(operationSystem)
	NetType     string = "net_type"     // NetType 连接网络 (netType)
	Language    string = "lang"         // Language 语言 (language)
	Locale      string = "locale"       // Locale 地区(locale)
)

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

func ApiOk(data interface{}) *ApiResult {
	return &ApiResult{
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
