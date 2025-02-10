package main

import "github.com/ywengineer/smart/utility"

type Configuration struct {
	Port             int                   `json:"port" yaml:"port"`
	BasePath         string                `json:"base_path" yaml:"base-path"`
	RDB              utility.RdbProperties `json:"rdb" yaml:"rdb"`
	Redis            string                `yaml:"redis" json:"redis"` // redis://user:password@host:port/?db=0&node=host:port&node=host:port
	RedisLock        bool                  `json:"redis_lock" yaml:"redis-lock"`
	MaxRequestBodyKB int                   `json:"max_request_body_kb,omitempty" yaml:"max-request-body-kb,omitempty"`
}
