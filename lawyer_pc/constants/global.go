package constants

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

var (
	API_HOST string = "https://admin.2kuan.com"
	Context *gin.Context = nil
)


