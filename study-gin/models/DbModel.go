package model

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
	dbreader *gorm.DB
	err      error
)

func InitDb() {
	// 这里这个配置将来要改成从Config读写类中获取
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		"root",
		"root",
		"127.0.0.1",
		"3306",
		"fastadmin",
		"utf8mb4",
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Println("数据库联接失败")
		os.Exit(1)
	}

	dsnreader := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		"reader",
		"llf",
		"127.0.0.1",
		"3308",
		"blog",
		"utf8mb4",
	)

	dbreader, err = gorm.Open(mysql.Open(dsnreader), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

}
