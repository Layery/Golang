package routers

import (
	v1 "study-gin/controllers/v1"
	"study-gin/middleware"

	"github.com/gin-gonic/gin"
)

func InitAdminRouter() {
	// 创建路由分组, 首先获取到框架路由
	router := gin.Default()
	router.Use(middleware.Log())
	adminRouter := router.Group("/v1")
	{
		// 在v1分组下注册路由
		//r.POST("/login_by_params", v1.LoginAction)
		adminRouter.POST("/login_by_struct", v1.LoginByStruct)
		adminRouter.POST("/struct", v1.LoadParamsByStruct)
		adminRouter.GET("/hello/:name", v1.SayHelloAction) // 可以通过:params的方式来接收url方式的get传参

		// 首页创建用户接口
		adminRouter.GET("/index", v1.IndexAction)

	}

	/**
	前台相关接口
	*/
	frontRouter := router.Group("/front")
	{
		frontRouter.GET("/test", v1.IndexController.GetDistList)
		frontRouter.GET("/files", v1.IndexController.TestFileReader)
		frontRouter.GET("/closure", v1.IndexController.TestClosure) // 练习闭包
		frontRouter.GET("/users", v1.UserController.GetUserList)    // 调用UserController下的一个方法
		frontRouter.GET("/map", v1.TestMapSlice)                    // 调用普通函数, 虽然也是UserController下的函数, 但是并不能直观的看出来, 它是在哪个文件下
		frontRouter.GET("/arr", v1.TestArrayMap)
	}

	_ = router.Run()
}
