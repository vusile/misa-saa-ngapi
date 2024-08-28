package model

import "gorm.io/gorm"

type Huduma struct {
	gorm.Model
	ID   uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name string
	SpecifyDays bool
}
