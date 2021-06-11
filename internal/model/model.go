package model

import (
	"fmt"
	"gin-api/common/global"
	"gin-api/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

func Init() {
	var err error
	//关于读写分离我们可以定义两个变量，一个读和一个写的，然后分别初始化，然后在查询和写入操作的时候使用不同的连接，或者使用一个map保存读和写的连接
	global.DBEngine, err = NewDBEngine()
	if err != nil {
		panic(err)
	}
}

func NewDBEngine() (*gorm.DB, error) {
	var dsn string
	cfg := config.Conf.Database
	switch cfg.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s",
			cfg.User,
			cfg.Password,
			cfg.Protocol,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)
	default:
		log.Fatalf("invalid db driver %v\n", cfg.Driver)
	}

	db, err := gorm.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatalf("Open "+cfg.Driver+" failed. %v\n", err)
		return nil, err
	}
	db.DB().SetConnMaxLifetime(cfg.MaxLifetime) //最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(cfg.MaxOpens)       //设置最大连接数
	db.DB().SetMaxIdleConns(cfg.MaxIdles)       //设置闲置连接数
	if cfg.RunMode == "debug" {
		db.LogMode(true)
	}
	//设置全局表名禁用复数
	db.SingularTable(true)
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return cfg.Prefix + defaultTableName
	//}
	//对updatedAt和createdAt字段自动更新
	//db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	//otgorm.AddGormCallbacks(db)
	return db, nil
}
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		//nowTime := time.Now().Unix()
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				_ = createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				_ = modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("UpdatedAt", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		isDelField, hasIsDelField := scope.FieldByName("IsDel")
		if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
			//now := time.Now().Unix()
			now := time.Now().Format("2006-01-02 15:04:05")
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v,%v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(now),
				scope.Quote(isDelField.DBName),
				scope.AddToVars(1),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
