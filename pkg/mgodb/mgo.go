package mgodb

import (
	"context"
	"gin-api/common/global"
	"gin-api/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"
)

//filter: bson.D{{"name", "aaa"}}
//update: bson.D{{"$set",
//		bson.D{
//			{"name", "aaa"},
//		},
//	}}
//https://mongoing.com/archives/27257 mongodb中文社区
type mgo struct {
	db     string
	col    string
	filter interface{}
	update interface{}
}

func Init() {
	cfg := config.Conf.Mongo
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Host).SetMaxPoolSize(cfg.MaxPoolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	global.Mongo = client
	err = global.Mongo.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
}

func Table(database, collection string) *mgo {
	return &mgo{
		db:  database,
		col: collection,
	}
}

// FindOne 查询单个
func (m *mgo) FindOne(key string, value interface{}) *mongo.SingleResult {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	//collection.
	singleResult := collection.FindOne(context.TODO(), bson.D{{key, value}})
}

//InsertOne 插入单个
func (m *mgo) InsertOne(value interface{}) *mongo.InsertOneResult {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	insertResult, err := collection.InsertOne(context.TODO(), value)
	if err != nil {
		log.Println(err)
		return nil
	}
	return insertResult
}

//CollectionCount 查询集合里有多少数据
func (m *mgo) CollectionCount() (string, int64) {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	name := collection.Name()
	size, _ := collection.EstimatedDocumentCount(context.TODO())
	return name, size
}

//CollectionDocuments 按选项查询集合 Skip 跳过 Limit 读取数量 sort 1 ，-1 . 1 为最初时间读取 ， -1 为最新时间读取
func (m *mgo) CollectionDocuments(Skip, Limit int64, sort int) *mongo.Cursor {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	SORT := bson.D{{"_id", sort}} //filter := bson.D{{key,value}}
	//filter := bson.D{{}}
	findOptions := options.Find().SetSort(SORT).SetLimit(Limit).SetSkip(Skip)
	//findOptions.SetLimit(i)
	temp, _ := collection.Find(context.Background(), m.filter, findOptions)
	return temp
}

//ParsingId 获取集合创建时间和编号
func (m *mgo) ParsingId(result string) (time.Time, uint64) {
	temp1 := result[:8]
	timestamp, _ := strconv.ParseInt(temp1, 16, 64)
	dateTime := time.Unix(timestamp, 0) //这是截获情报时间 时间格式 2019-04-24 09:23:39 +0800 CST
	temp2 := result[18:]
	count, _ := strconv.ParseUint(temp2, 16, 64) //截获情报的编号
	return dateTime, count
}

//DeleteAndFind 删除和查询
func (m *mgo) DeleteAndFind(key string, value interface{}) (int64, *mongo.SingleResult) {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	singleResult := collection.FindOne(context.TODO(), m.filter)
	DeleteResult, err := collection.DeleteOne(context.TODO(), m.filter, nil)
	if err != nil {
		log.Println("删除时出现错误，你删不掉的~")
		return 0, nil
	}
	return DeleteResult.DeletedCount, singleResult
}

//Delete 删除
func (m *mgo) Delete(key string, value interface{}) int64 {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	count, err := collection.DeleteOne(context.TODO(), m.filter, nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	return count.DeletedCount

}

//DeleteMany 删除多个
func (m *mgo) DeleteMany(key string, value interface{}) int64 {
	collection := global.Mongo.Database(m.db).Collection(m.col)
	count, err := collection.DeleteMany(context.TODO(), m.filter)
	if err != nil {
		log.Println(err)
		return 0
	}
	return count.DeletedCount
}

func (m *mgo) FindAndUpdate(filter interface{}, update interface{}) *mongo.SingleResult {
	return global.Mongo.Database(m.db).Collection(m.col).FindOneAndUpdate(context.TODO(), m.filter, m.update)
}

func (m *mgo) UpdateOne(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return global.Mongo.Database(m.db).Collection(m.col).UpdateOne(context.TODO(), m.filter, m.update)
}

func (m *mgo) UpdateMany(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return global.Mongo.Database(m.db).Collection(m.col).UpdateMany(context.TODO(), m.filter, m.update)
}
