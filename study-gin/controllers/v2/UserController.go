package v2

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// v2版本的登录接口
func LoginAction(c *gin.Context) {
	// 获取Get 请求参数
	username := c.Query("username") // query方法, 可以获取get参数,

	username1 := c.DefaultQuery("username", "weidingyi") // defaultQuery方法, 同样可以获取get参数, 并支持设置一个默认值

	username2, isset := c.GetQuery("username") // 该方法有两个返回值, 第二个返回值返回给定的key在get参数中是否存在, 类似于PHP的isset函数的用法


	c.String(200, "username %s,  username1 %s, username2 %s, isset %v", username, username1, username2, isset)

}

func GetUserListAction() {
	fmt.Println("hello getUserListAction")
}
