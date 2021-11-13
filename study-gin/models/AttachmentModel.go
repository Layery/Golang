package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type Attachment struct {
	gorm.Model
	gorm.DeletedAt  `gorm:"column:createtime"`
}

func (Attachment) TableName() string {
	return "fa_attachement"
}

func (a *Attachment) GetList(list []Attachment) interface{} {

	//log.Printf("list is : %#v", list)
	body, _ := json.MarshalIndent(list, "", "\t")

	log.Printf("%s \n", body)
	//res := DB.Find(&list)
	//return res
	return nil
}


