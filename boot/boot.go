package boot

import (
	"gin-api/config"
	"gin-api/internal/model"
	"gin-api/pkg/gredis"
	"gin-api/pkg/logger"
	"gin-api/pkg/redigo"
)

func init() {
	config.Init()
	logger.Init()
	model.Init()
	redigo.Init() //使用redisgo操作redis，和下面二选一
	gredis.Init() //使用go-redis操作redis
	//tracer.Init()
	//es.Init()
	//mgodb.Init()
}
