package middleware

import (
	"ginpro/common/global"
	"ginpro/pkg/app"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(map[string]interface{})
		data["ip"] = c.ClientIP()
		data["url"] = c.Request.Host + c.Request.RequestURI
		data["method"] = c.Request.Method
		data["proto"] = c.Request.Proto
		data["request"] = app.JsonParams(c)
		data["header"] = c.Request.Header

		c.Next()

		data["response"] = app.Get(c, "responseData")
		// 写日志
		//level := app.GetLevel(c)
		level := app.GetLevel(c)
		switch level {
		case "error":
			global.Logger.WithFields(data).Error(app.GetDetail(c))
			break
		case "warn":
			global.Logger.WithFields(data).Warn()
			break
		default:
			global.Logger.WithFields(data).Info()
			break
		}
	}
}
