package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	DBEngine  *gorm.DB
	RedisPool *redis.Pool
	Logger    *logrus.Logger
	Tracer opentracing.Tracer
	Es *elasticsearch.Client
	Mongo *mongo.Client
)
var (
	//Env 运行环境
	Env string
	//RequestID 请求ID
	RequestID int64
	//RequestTime 请求开始时间
	RequestTime time.Time
	//RequestData 请求参数
	RequestData interface{}
	//Ip 客户端IP
	Ip string
	//LinkMax 友链最大数
	LinkMax int
	//PageSize 默认分页数
	PageSize int
	//LabelMax 标签最大数
	LabelMax int
)
