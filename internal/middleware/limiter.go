package middleware

import (
	"ginpro/common/dict"
	"ginpro/pkg/app"
	"ginpro/pkg/limiter"
	"github.com/gin-gonic/gin"
)

func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if bucket, ok := l.GetBucket(key); ok {
			count := bucket.TakeAvailable(1)
			if count == 0 {
				app.Error(c, dict.TooManyRequests)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
