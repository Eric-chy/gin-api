package middleware

import (
	"gin-api/common/dict"
	"gin-api/pkg/app"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = dict.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}
		if token == "" {
			ecode = dict.InvalidParams
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = dict.UnauthorizedTokenTimeout
				default:
					ecode = dict.UnauthorizedTokenError
				}
			}
		}

		if ecode != dict.Success {
			app.Error(c, ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
