package middleware

import (
	"fmt"
	"gin-api/common/dict"
	"gin-api/pkg/app"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Sprintf("错误信息:：%s", err))
				app.Error(c, dict.ServerError.WithDetails(fmt.Sprintf("错误信息:：%s", err)).WithLevel("error"))
				c.Abort()
			}
		}()

		c.Next()
	}
}
