package model

import (
	"gorm.io/gorm"
)

type Parokia struct {
	gorm.Model
	ID        uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name      string
	JimboID   uint
	Jimbo     Jimbo
	IsKigango bool
}
