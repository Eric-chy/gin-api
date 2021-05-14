package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	rdb "github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DBEngine  *gorm.DB
	RedisPool *redis.Pool
	Redis     *rdb.Client
	Logger    *logrus.Logger
	Tracer    opentracing.Tracer
	Es        *elasticsearch.Client
	Mongo     *mongo.Client
)
