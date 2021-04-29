package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Global() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置接口请求开始时间
		c.Set("beginTime", time.Now())
		c.Next()
	}
}
