package main

import (
	"log"
	"skill-order/server/dao"
)

func init() {
	// 链接初始化
	dao.InitMysql()

	if dao.DB == nil {
		log.Fatal("db 初始化失败!")
	}

	// 商品初始化
	dao.InitGoodsData()

}

func main() {

	// 检查redis链接, rabbitmq链接, mysql链接

	// 启动server层服务
	// server.Start()

	//r := InitRoute()

	//_ = r.Run(":9000")

}
