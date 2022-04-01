package main

import (
	"fmt"
	"lawyerpc/constants"
	"lawyerpc/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db  *gorm.DB
	err error
)

func main() {
	gin.SetMode("debug")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		"divorce",
		"GzEi6TEPJ5n8GrWa",
		"39.105.27.31",
		"3306",
		"divorce",
		"utf8mb4",
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("database connect failed")
		os.Exit(1)
	} else {
		constants.DB = db
		log.Printf("database is %v", db)
	}
	r := router.InitFrontEndRouter()
	err = r.Run(":8088")
	log.Printf("System Error is %v", err)

}
