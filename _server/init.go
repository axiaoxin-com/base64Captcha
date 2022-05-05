// 依赖服务初始化

package main

import (
	"fmt"
	"time"

	"github.com/axiaoxin-com/base64Captcha"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

const (
	// AppID 服务app id
	AppID = 5
)

var (
	// RedisClient redis 客户端
	RedisClient *redis.Client
	// RedisStore redis store
	RedisStore base64Captcha.Store
)

// InitLogging logging初始化
func InitLogging() {
	level := viper.GetString("logging.level")
	if level != "" {
		logging.SetLevel(level)
	}
}

// InitRedis 初始化Redis
func InitRedis() {
	if RedisClient == nil {
		cli, err := goutils.RedisClient(fmt.Sprintf("%s", viper.GetString("env")))
		if err != nil {
			logging.Fatal(nil, "init redis error: "+err.Error())
		}
		RedisClient = cli
		RedisStore = base64Captcha.NewRedisStore(cli, viper.GetDuration("captcha.expire_seconds")*time.Second, "vcode")
	}
}
