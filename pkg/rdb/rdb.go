package rdb

import (
	"fmt"
	"github.com/go-gorm/caches/v4"
	"github.com/pkg/errors"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
	"strings"
	"time"
)

// Properties rational database configuration properties
type Properties struct {
	Name       string           `json:"name" yaml:"name"` // mysql or postgres
	Username   string           `json:"username" yaml:"username"`
	Password   string           `json:"password" yaml:"password"`
	Host       string           `json:"host" yaml:"host"`
	Port       int              `json:"port" yaml:"port"`
	Database   string           `json:"database" yaml:"database"`
	Parameters string           `json:"parameters" yaml:"parameters"`
	Pool       DbPoolProperties `json:"pool" yaml:"pool"`
	Cache      string           `json:"cache" yaml:"cache"`
	DebugMode  bool             `json:"debug_mode" yaml:"debug-mode"`
}

type DbPoolProperties struct {
	MaxIdleCon          int   `json:"max_idle_con" yaml:"max-idle-con"`
	MaxOpenCon          int   `json:"max_open_con" yaml:"max-open-con"`
	MaxLifeTimeInMinute int64 `json:"max_life_time_minute" yaml:"max-life-time-minute"`
}

// NewRDB create rational database instance
func NewRDB(driver utility.RdbProperties, plugins ...gorm.Plugin) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var cachePlugin gorm.Plugin
	//
	if driver.Name == "mysql" {
		db, err = NewMySQL(driver)
	} else if driver.Name == "postgres" {
		db, err = NewPostgres(driver)
	} else {
		db, err = nil, errors.New("not support driver: "+driver.Name)
	}
	if err != nil {
		return nil, err
	} else if len(driver.Cache) > 0 { // cache
		if strings.HasPrefix(driver.Cache, "mem://") {
			if memProtocol, err := url.Parse(driver.Cache); err == nil {
				cachePlugin = &caches.Caches{Conf: &caches.Config{
					Cacher: (&memoryCacher{}).size(utility.QueryInt(memProtocol.Query(), "size")),
				}}
			} else {
				utility.DefaultLogger().Error("rdb cache inactivate, because of create failed: " + driver.Cache)
			}
		} else if strings.HasPrefix(driver.Cache, "redis://") {
			cachePlugin = &caches.Caches{Conf: &caches.Config{
				Cacher: &redisCacher{rdb: utility.NewRedis(driver.Cache)},
			}}
		} else {
			utility.DefaultLogger().Error("rdb not support this cache: " + driver.Cache)
		}
	}
	//
	if db != nil {
		if cachePlugin != nil { // // cache plugin
			_ = db.Use(cachePlugin)
		}
		if plugins != nil && len(plugins) > 0 {
			for _, plugin := range plugins {
				_ = db.Use(plugin)
			}
		}
	}
	//
	return initRbdConnPool(db, driver)
}

// NewMySQL create gorm.DB instance based on mysql database
func NewMySQL(driver utility.RdbProperties) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", driver.Username, driver.Password, driver.Host, driver.Port, driver.Database)
	if len(driver.Parameters) > 0 {
		dsn += "&" + driver.Parameters
	}
	//
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // autoconfigure based on currently MySQL version
	}), defaultConfig(driver.DebugMode))
}

// NewPostgres create gorm.DB instance based on postgres database
func NewPostgres(driver utility.RdbProperties) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=Asia/Shanghai",
		driver.Host, driver.Port, driver.Username, driver.Password, driver.Database)
	if len(driver.Parameters) > 0 {
		dsn += " " + driver.Parameters
	}
	//
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), defaultConfig(driver.DebugMode))
}

func initRbdConnPool(gdb *gorm.DB, driver utility.RdbProperties) (*gorm.DB, error) {
	db, err := gdb.DB()
	if err != nil {
		utility.DefaultLogger().Error("get db instance from gorm error", zap.Any("driver", driver), zap.Error(err))
		return nil, err
	}
	db.SetMaxIdleConns(utility.MaxInt(driver.Pool.MaxIdleCon, 5))
	db.SetMaxOpenConns(utility.MaxInt(driver.Pool.MaxOpenCon, 5))
	db.SetConnMaxLifetime(time.Duration(utility.MaxInt64(1, driver.Pool.MaxLifeTimeInMinute) * int64(time.Minute)))
	if err = db.Ping(); err != nil {
		utility.DefaultLogger().Error("connect to db instance failed", zap.Any("driver", driver), zap.Error(err))
		return nil, err
	}
	return gdb, nil
}

func defaultConfig(debug bool) *gorm.Config {
	if debug {
		return &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Info),
		}
	}
	return &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}
}
