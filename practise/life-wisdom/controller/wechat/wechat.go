package wechat

import (
	"wisdom/controller"

	"github.com/gin-gonic/gin"
)

type Wechat struct {
	controller.Base
}

func NewWechat() *Wechat {
	return &Wechat{}
}

func (wechat *Wechat) ListImage(c *gin.Context) {
	c.JSON(200, gin.H{"data": "list image"})
}
