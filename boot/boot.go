package boot

import (
	"ginpro/config"
	"ginpro/internal/model"
	"ginpro/pkg/gredis"
	"ginpro/pkg/logger"
	"ginpro/pkg/redigo"
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
