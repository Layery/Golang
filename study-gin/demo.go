package main

import (
	"github.com/gin-gonic/gin"
	"study-gin/models"
	"study-gin/routers"
)

func main() {
	// 初始化一个http服务对象
	//router := gin.Default()
	//fmt.Printf("%v", router)

	// gin框架默认工作在debug模式下
	// 可以通过代码设置工作模式
	gin.SetMode("debug") // debug, release, test


	// 初始化db联接
	model.InitDb()


	// 初始化后台路由
	routers.InitAdminRouter()

	// 初始化前台路由 <--- 为啥你不起作用???
	//routers.InitFrontRouter()

}
