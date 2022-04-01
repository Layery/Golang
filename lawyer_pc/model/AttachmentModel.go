package model

import (
	"errors"
	"fmt"
	"lawyerpc/constants"
)

type AttachmentModel struct {
	//BaseModel
	//ID  int    `json:"id"`
	URL string `json:"url"`
}

func (AttachmentModel) TableName() string {
	return "fa_attachment"
}

func (att *AttachmentModel) GetAttachList(category string, isSingle bool) (interface{}, error) {

	isExists := true
	constants.DB.Table(att.TableName()).Select("exists(select * from fa_attachment where category = '"+ category+"')").Take(&isExists)
	if !isExists {
		return []int{}, errors.New("the category can not found")
	}

	if isSingle == true {
		return att.getOne(category), nil
	} else {
		var attachment []AttachmentModel // 声明一个切片但并不赋值
		constants.DB.Where(map[string]interface{}{
			"category": category,
		}).Find(&attachment)
		return attachment, nil
	}
}

func (att *AttachmentModel) getOne(category string) interface{} {
	var model AttachmentModel
	constants.DB.Where(map[string]interface{}{
		"category": category,
	}).Take(&model)
	result := map[string]interface{}{
		category: model.URL,
	}
	return result
}

func (att *AttachmentModel) GetAttachmentLink(item *ArticleModel) {
	fmt.Printf("came in params type is: %T, value is : %#v \n\n", item, item)
	item.TitleImg = "www" + item.TitleImg
}

func GetAttachmentPath(item *ArticleModel){
	item.TitleImg = constants.API_HOST + item.TitleImg
}































