package boot

import (
	"ginpro/config"
	"ginpro/pkg/gredis"
	"ginpro/pkg/logger"
)

func init() {
	config.Init()
	logger.Init()
	//model.Init()
	gredis.Init()
	//tracer.Init()
	//es.Init()
	//mgodb.Init()
}
