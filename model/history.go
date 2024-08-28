package model

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	ID          uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	UserID      uint
	ModelID     uint
	ModelTypeID uint
	From        *string `gorm:"type:text"`
	To          string  `gorm:"type:text"`
}
