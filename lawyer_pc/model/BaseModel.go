package model

import "gorm.io/gorm"

type BaseModel struct {
	*gorm.DB
}
