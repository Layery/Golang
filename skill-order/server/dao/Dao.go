package dao

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"skill-order/server/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// mysql 的配置, 这个将来改成走配置文件
type mysqlConfig struct {
	user     string
	password string
	host     string
	port     int16
	dbname   string
}

var cfg = mysqlConfig{
	user:     "root",
	password: "root",
	host:     "127.0.0.1",
	port:     3306,
	dbname:   "skill_study",
}

// 初始化mysql
func InitMysql() (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.user, cfg.password, cfg.host, cfg.port, cfg.dbname)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info), // 设置日志级别为 Info
	})
	if err != nil {
		fmt.Println("Failed to gorm open dsn err:", err)
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get database connection:", err)
		return
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	return err
}

func InitGoodsData() {
	if DB.Migrator().HasTable(&model.Order{}) { // 判断表是否存在
		DB.Migrator().DropTable(&model.SkillGoods{})
		DB.Migrator().DropTable(&model.Order{})
	}

	_ = DB.AutoMigrate(&model.SkillGoods{}, &model.Order{})

	// 给商品表来一点数据
	var cnt int64
	DB.Model(&model.SkillGoods{}).Count(&cnt)

	if cnt <= 0 {
		goodsList := []model.SkillGoods{}
		for i := 0; i < 10; i++ {
			goodsList = append(goodsList, model.SkillGoods{
				Name:  "商品" + fmt.Sprintf("%d", i+1),
				Price: 100.05,
				Stock: 100,
			})
		}
		DB.Create(goodsList)
	}
}
