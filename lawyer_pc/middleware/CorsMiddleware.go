package middleware

import (
	"github.com/gin-gonic/gin"
	"lawyerpc/constants"
	"net/http"
)

func CorsNext() gin.HandlerFunc {
	return func(c *gin.Context) {
		constants.Context = c
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		//c.Header("Access-Control-Allow-Origin", "http://test.cc")
		c.Header("x-run-with-go", "yes")
		c.Next()
	}
}
