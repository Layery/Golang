package model

import (
	"errors"
	"fmt"
	"lawyerpc/constants"
	"lawyerpc/utils"
)

type ArticleModel struct {
	ID          int    `json:"id"`
	CateID      int    `json:"cate_id"`
	Title       string `json:"title"`
	ContentDesc string `json:"content_desc"`
	Content     string `json:"content"`
	TitleImg    string `json:"title_img"`
	Createtime  string `json:"createtime"`
	Updatetime  string `json:"updatetime"`
}

func (ArticleModel) TableName() string {
	return "fa_cms"
}

func (article ArticleModel) GetCateList() (interface{}, error) {
	type CateModel struct {
		Id int `json:"id"`
		Name string `json:"name"`
	}
	var datalist []CateModel
	res := constants.DB.Table("fa_category").Find(&datalist)
	if res.RowsAffected <= 0 {
		return []int{}, errors.New("未找到栏目列表")
	}
	return datalist, nil
}

func (article *ArticleModel) GetList(page, cateID int) (interface{}, error) {
	var (
		articleSlice []ArticleModel
		total        int64
	)
	if cateID <= 0 {
		var categoryModel struct {
			Id   int
			Name string
		}
		res := constants.DB.Table("fa_category").Select([]string{
			"id", "name",
		}).Where(map[string]interface{}{
			"type": "article",
		}).Take(&categoryModel)
		if res.RowsAffected <= 0 {
			datalist := map[string]interface{}{
				"total": total,
				"data":  articleSlice,
			}
			return datalist, errors.New("参数有误!")
		}
		cateID = categoryModel.Id
	}
	res := constants.DB.Select([]string{
		"id",
		"cate_id", // id, cate_id, title, content_desc, title_img, createtime, updatetime
		"title",
		"content_desc",
		"title_img",
		"createtime",
		"updatetime",
	}).Where(map[string]interface{}{
		"cate_id": cateID,
	}).Limit(10).Offset(10 * (page - 1)).Order("id desc").Find(&articleSlice)
	constants.DB.Table(article.TableName()).Where(map[string]interface{}{"cate_id": cateID}).Count(&total)

	for key := range articleSlice {
		GetAttachmentPath(&articleSlice[key])
		articleSlice[key].Createtime = utils.FormatTime(articleSlice[key].Createtime)
		articleSlice[key].Updatetime = utils.FormatTime(articleSlice[key].Updatetime)
	}

	datalist := map[string]interface{}{
		"total": total,
		"data":  articleSlice,
	}
	if res.RowsAffected <= 0 {
		return datalist, errors.New("not found")
	}
	return datalist, nil

}

func (article *ArticleModel) GetArticleDetail(articleID int) (interface{}, error)  {
	fmt.Printf("=====> %#v \n", articleID)
	if articleID <= 0 {
		return []int{}, errors.New("文章id不得为空")
	}
	var data ArticleModel
	constants.DB.Where(map[string]interface{}{
		"id" : articleID,
	}).Take(&data)

	fmt.Printf("before data is %#v \n", data)

	data.Createtime = utils.FormatTime(data.Createtime)
	data.Updatetime = utils.FormatTime(data.Updatetime)
	data.TitleImg = constants.API_HOST + data.TitleImg

	fmt.Printf("data is %#v \n", data)
	return data, nil
}