package controller

import (
	"fmt"
	"io/ioutil"
	"lawyerpc/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseController
}

var Index = new(IndexController)

func (index *IndexController) GetSystemInfo(c *gin.Context) {
	res, _ := http.Get(constants.API_HOST + "/api/index/systeminfo")
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(http.StatusOK, fmt.Sprintf("%s", body))
}
