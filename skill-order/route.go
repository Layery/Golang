package main

import (
	"github.com/gin-gonic/gin"
	"skill-order/server"
)

func InitRoute() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		// 使用for update
		v1.GET("/skill", server.NewOrderController().SkillV1)
	}

	v2 := r.Group("/v2")
	{
		// 使用乐观锁, 用update_at作为版本号
		v2.GET("/skill", server.NewOrderController().SkillV2)
	}

	return r
}
