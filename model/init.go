package model

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var _db *gorm.DB

func Database(conn string) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,
		DefaultStringSize:         256,  //string默认长度
		DisableDatetimePrecision:  true, //禁止datetime精度，mysql5.6之前不支持
		DontSupportRenameIndex:    true, //重命名索引，就要把索引先删除再重建，mysql5.7之前不支持
		DontSupportRenameColumn:   true, //用change重命名列，mysql8之前的数据库不支持
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //不在表后加s
		},
	})
	if err != nil {
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(20)  //设置连接池
	sqlDB.SetConnMaxIdleTime(100) //打开连接池
	sqlDB.SetConnMaxIdleTime(time.Second * 30)

	_db = db
	log.Println("mysql connect successfully")
	migration()

}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
