package middleware

import (
	"gin-api/common/dict"
	"gin-api/pkg/app"
	"gin-api/pkg/limiter"
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
