package middleware

import (
	"ginpro/common/global"
	"ginpro/pkg/app"
	"ginpro/pkg/helper/gtime"
	"github.com/gin-gonic/gin"
	"time"
)

func Global() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置接口请求开始时间
		global.RequestTime = time.Now()
		// 设置请求ID
		global.RequestID = gtime.GetMicroTime()
		// 设置请求参数
		global.RequestData = app.JsonParams(c)
		// 设置客户端ip
		global.Ip = c.ClientIP()
		c.Next()
	}
}
