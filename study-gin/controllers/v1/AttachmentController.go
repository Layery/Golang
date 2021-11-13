package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	model "study-gin/models"
)
type AttachmentControllerStruct struct {
	BaseController
}

var AttachmentController = new(AttachmentControllerStruct)


func (u *AttachmentControllerStruct) GetList(c *gin.Context) {

	list := make([]model.Attachment, 5)

	//c.JSON(http.StatusOK, gin.H{
	//	"list": list,
	//})
	//rs := model.DB.Table("fa_attachment").Select("id").Find(&list)

	model := model.Attachment{}

	rs := model.GetList(list)

	c.JSON(http.StatusOK, gin.H{
		"list": rs,
	})
}


