package controller

import (
	"lawyerpc/model"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	BaseController
}

var Article = new(ArticleController)

func (article *ArticleController) GetArticleList(c *gin.Context) {
	articleModel := model.ArticleModel{}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	type TempPostStruct struct {
		CateId int `json:"cate_id"`
		Page   int `json:"page"`
	}
	var post TempPostStruct
	_ = c.BindJSON(&post)
	datalist, _ := articleModel.GetList(page, post.CateId)
	article.SuccessReturn("success", datalist, c)
}

func (article *ArticleController) GetCateList(c *gin.Context) {
	datalist, _ := (new(model.ArticleModel)).GetCateList()
	article.SuccessReturn("success", datalist, c)
}

func (article *ArticleController) GetDetail(c *gin.Context)  {
	type TempPostStruct struct {
		ArticleId string `json:"article_id"`
	}
	var post TempPostStruct
	err := c.ShouldBindJSON(&post)
	if err != nil {
		log.Printf("error is : %#v\n ", err.Error())
	}
	articleModel := model.ArticleModel{}
	articleID, _ := strconv.Atoi(post.ArticleId)
	data, _ := articleModel.GetArticleDetail(articleID)
	article.SuccessReturn("success", data, c)
}