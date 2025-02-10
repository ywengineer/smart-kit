package pkg

const (
	ContextKeySmart   = "smart-context"
	ContextKeyRDB     = "rdb"
	ContextKeyRedis   = "redis"
	ContextKeyApp     = "app"
	ContextKeyRedLock = "red-lock"
)

type ApiResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
