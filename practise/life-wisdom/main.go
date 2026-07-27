package main

import (
	"wisdom/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	initLogger()
	defer logger.Sync()

	// 采用gin框架, 监听微信输入请求, 当请求发来后, 自动调动生图接口, 自动上传微信公众号后台,
	r := gin.New()

	// 初始化路由
	r = router.InitRouter(r)

	_ = r.Run(":8080")
}

func initLogger() {
	obj, _ := zap.NewProduction()
	logger = obj.Sugar()
}
