package model

import (
	"gorm.io/gorm"
)

type Jimbo struct {
	gorm.Model
	ID         uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name       string `gorm:"uniqueIndex;size:40"`
	IsJimboKuu bool
	CountryID  uint
	Country    Country
	ChurchID   uint
	Church     Church
	Slug       string
}
