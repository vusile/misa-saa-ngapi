package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name         string `gorm:"type:varchar(50)"`
	Phone        string `gorm:"uniqueIndex;type:varchar(20)"`
	ActivatedAt  *time.Time
	Code         int
	SessionToken string `gorm:"index;type:varchar(100)"`
	CsrfToken    string `gorm:"type:varchar(100)"`
	Password     string `gorm:"type:varchar(100)"`
	ChurchID     uint
	CountryID    uint
}
