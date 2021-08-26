package v1

import "github.com/gin-gonic/gin"

type BaseController struct {
	ModuleName string
	ControllerName string
	ActionName string
}

// 基类下一个公共的json返回方法
func (this *BaseController) JsonReturn(code int, data interface{}, c *gin.Context) {
	c.JSON(code, gin.H{
		"data": data,
	})
}



