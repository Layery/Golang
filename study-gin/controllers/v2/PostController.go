package v2

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	c.String(500, "hello v2")
}

func est() {
	fmt.Println("hello getUserListAction")
}
