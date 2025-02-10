package pkg

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SmartContext struct {
	Rdb     *gorm.DB
	Redis   redis.UniversalClient
	RedLock *redislock.Client
}
