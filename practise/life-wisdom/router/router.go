package router

import (
	"wisdom/controller/image"
	"wisdom/controller/wechat"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) *gin.Engine {
	// 生成图片组
	imageRouter := router.Group("/img")
	{
		imageController := image.NewImage()
		commonImageController := image.NewCommonImage()

		imageRouter.GET("/genImage", imageController.GenImage)
		imageRouter.GET("/common", commonImageController.GenCommonImage)
	}

	// 微信相关组
	wechatRouter := router.Group("/wx")
	{
		wechatController := wechat.NewWechat()

		wechatRouter.GET("/list", wechatController.ListImage)

	}
	return router
}
