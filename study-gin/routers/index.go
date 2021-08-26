package routers

import (
	"github.com/gin-gonic/gin"
	"log"
	"study-gin/controllers/v1"
)

func InitFrontRouter() {
	router := gin.Default()
	frontRouter := router.Group("/v1")
	{
		// 再v1下注册路由
		log.Println("index.go register router")
		frontRouter.GET("/users", v1.UserController.GetUserList)
	}

	_ = router.Run()
}