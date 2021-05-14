package middleware

import (
	"ginpro/common/dict"
	"ginpro/common/global"
	"ginpro/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"strconv"
	"strings"
	"time"
)

func RedisLimiter(rate int, t ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rdb := global.Redis
		limiter := redis_rate.NewLimiter(rdb)
		uri := c.Request.RequestURI
		index := strings.Index(uri, "?")
		var key string
		if index == -1 {
			key = uri
		} else {
			key = uri[:index]
		}
		var res *redis_rate.Result
		var err error
		if len(t) > 0 && t[0] == "s" {
			res, err = limiter.Allow(c, key, redis_rate.PerSecond(rate))
		} else {
			res, err = limiter.Allow(c, key, redis_rate.PerMinute(rate))
		}
		if err != nil {
			global.Logger.Error(err)
		} else {
			c.Header("RateLimit-Remaining", strconv.Itoa(res.Remaining))
			if res.Allowed == 0 {
				seconds := int(res.RetryAfter / time.Second)
				c.Header("RateLimit-RetryAfter", strconv.Itoa(seconds))
				app.Error(c, dict.ErrRateLimited)
				c.Abort()
			}
		}
		c.Next()
	}
}
