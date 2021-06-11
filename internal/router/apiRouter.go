package router

import (
	"gin-api/config"
	_ "gin-api/docs"
	"gin-api/internal/api"
	"gin-api/internal/middleware"
	"gin-api/pkg/limiter"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"time"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})
var rate = 2

func ApiRouter() *gin.Engine {
	var r *gin.Engine
	// 创建一个不包含中间件的路由器
	r = gin.New()
	//// 使用 Global 中间件
	r.Use(middleware.Global())
	// 使用 Cors 中间件
	r.Use(middleware.Cors())
	//// 使用 Logger 中间件
	r.Use(middleware.Logger())
	//// 使用 Recovery 中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RateLimiter(methodLimiters)) //单机限流
	r.Use(middleware.RedisLimiter(rate))          //分布式限流
	r.Use(middleware.ContextTimeout(60 * time.Second))
	r.Use(middleware.Translations())
	//r.Use(middleware.Tracing())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//url := ginSwagger.URL("http://127.0.0.1:8000/swagger/doc.json")
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	atl := r.Group("../api")
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	//文件访问
	r.StaticFS("/static", http.Dir(config.Conf.App.UploadSavePath))
	article1 := api.NewArticle()
	atl.GET("/articles", article1.ArticleList)
	atl.GET("/articles/:id", article1.ArticleDetail)
	//测试发邮件，需要先配置好
	atl.GET("/articles/email", article1.SendEmail)
	//测试httpclient，类似于curl
	atl.GET("/articles/curl", article1.Curl)
	//article := r.Group("articles")
	//{
	//	article.GET("", api.ArticleList)
	//	article.GET(":id", api.ArticleDetail)
	//}

	return r
}
