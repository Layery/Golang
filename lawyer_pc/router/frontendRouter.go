package router

import (
	"lawyerpc/controller"
	"lawyerpc/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitFrontEndRouter() *gin.Engine {
	//var router *gin.Engine

	router := gin.Default()

	router.Use(middleware.CorsNext())

	router.POST("", func(c *gin.Context) {
		c.String(http.StatusOK, "hello weidingyi go!")
	})

	// 文件服务器
	router.StaticFS("/files", http.Dir("D:/"))

	// api
	api := router.Group("api")
	{
		api.POST("/cms/articlelist", controller.Article.GetArticleList)
		api.POST("/cms/catelist", controller.Article.GetCateList)
		api.POST("/cms/articledetail", controller.Article.GetDetail)
		api.POST("/index/systeminfo", controller.Index.GetSystemInfo)
		api.POST("/index/get_img", controller.Base.GetImgByCategory)
	}
	return router
}
