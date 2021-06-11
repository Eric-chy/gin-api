package gredis

import (
	"context"
	"fmt"
	"gin-api/common/global"
	"gin-api/config"
	"github.com/go-redis/redis/v8"
	"time"
)

func Init() {
	cfg := config.Conf.Redis
	global.Redis = redis.NewClient(&redis.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		DialTimeout:  5 * time.Second, //不设置默认值也是5
		ReadTimeout:  3 * time.Second, //此处是默认值，也可以不设置或者配置文件里配置
		WriteTimeout: 3 * time.Second, //此处是默认值，也可以不设置
		PoolSize:     20,
	})
	if err := global.Redis.Ping(context.Background()).Err(); err != nil {
		fmt.Println("redis connect fail:", err)
		panic(err)
	}
}
