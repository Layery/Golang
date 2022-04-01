package controller

import (
	"github.com/gin-gonic/gin"
	"lawyerpc/model"
	"net/http"
)

type BaseController struct {
	ControllerName string
	ActionName string
}


var Base = new(BaseController)

// GetImgByCategory 根据类别获取图片的方法
func (Base *BaseController) GetImgByCategory(c *gin.Context)  {
	category := c.PostForm("category")
	if category == "" {
		Base.ErrorReturn("参数有误: category", []interface{}{}, c)
	}
	attach := new(model.AttachmentModel)
	res, err := attach.GetAttachList(category, false)
	if err != nil {
		Base.ErrorReturn(err.Error(), res, c)
	} else {
		Base.SuccessReturn("success", res, c)
	}
}

func (Base *BaseController) SuccessReturn(msg string, data interface{}, c *gin.Context)  {
	ajaxReturn(1, msg, data, c)
}
func (Base *BaseController) ErrorReturn(msg string, data interface{}, c *gin.Context)  {
	ajaxReturn(0, msg, data, c)
}

func ajaxReturn(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		//"data": map[string]interface{}{
		//	"name": "liulongfei",
		//	"list": list,
		//	"array": []map[string]interface{}{
		//		{"aaa": "aaa", "bbb": "bbb"},
		//		{"aaa": "aaa", "bbb": "bbb"},
		//	},
		//	"age": 30,
		//},
		"data": data,
	})
}

func (Base *BaseController) GetFileList(c *gin.Context)  {
	
}