package model

import (
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	ID   uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name string `gorm:"uniqueIndex;size:40"`
}
