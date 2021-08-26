package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Log() gin.HandlerFunc {
	return func(context *gin.Context) {
		ip := context.ClientIP()
		method := context.Request.Method
		url := context.FullPath()

		msg := fmt.Sprintf("当前的请求ip: %v, 请求方法: %v, 请求路径: %v \n\n", ip, method, url)

		log.Println(msg)
	}
}