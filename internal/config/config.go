package config

import (
	"sync"

	"github.com/spf13/viper"

	"github.com/redis/go-redis/v9"
	"github.com/xmopen/golib/pkg/xconfig"
	"github.com/xmopen/golib/pkg/xlogging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	configInstance *xconfig.Config

	blogsDataBaseInstance *gorm.DB
	blogsDataBaseOnce     sync.Once
	// blogsDataBaseDNS root:xxx@/openxm?charset=utf8&parseTime=True&loc=Local
	blogsDataBaseDNS string

	redisInstance *redis.Client
	redisOnce     sync.Once
	redisAddr     string
	redisPass     string

	xlog = xlogging.Tag("blogsvr.config")
)

func init() {
	configInstance = xconfig.InitConfig()
	blogsDataBaseDNS = configInstance.Config().GetString("database.mysql.dns")
	redisAddr = configInstance.Config().GetString("database.redis.addr")
	redisPass = configInstance.Config().GetString("database.redis.pass")
}

// Config return viper config
func Config() *viper.Viper {
	return configInstance.Config()
}

// AuthDataBase authDataBase.
func AuthDataBase() *gorm.DB {
	if blogsDataBaseInstance == nil {
		blogsDataBaseOnce.Do(func() {
			dbInstance, err := gorm.Open(mysql.Open(blogsDataBaseDNS), &gorm.Config{})
			if err != nil {
				xlog.Errorf("open mysql err:[%+v]", err)
				return
			}
			blogsDataBaseInstance = dbInstance
		})
	}
	return blogsDataBaseInstance
}

// GlobalRedis globalRedis实例.
func GlobalRedis() *redis.Client {
	if redisInstance == nil {
		redisOnce.Do(func() {
			redisInstance = redis.NewClient(&redis.Options{
				Addr:     redisAddr,
				Password: redisPass,
			})
		})
	}
	return redisInstance
}
