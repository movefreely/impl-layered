package sqlinit

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func SqlInit() *gorm.DB {
	dsn := "im:oushixing@tcp(121.4.86.52:3306)/im?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库出错")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接出错")
	}
	sqlDB.SetMaxIdleConns(15)
	sqlDB.SetMaxOpenConns(100)
	DB = db
	return db
}
