package model

import "gorm.io/gorm"

type WeekDay struct {
	gorm.Model
	ID       uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name     string `gorm:"uniqueIndex;size:10"`
	Priority int
}
