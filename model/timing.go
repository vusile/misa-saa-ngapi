package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Timing struct {
	gorm.Model
	ID         uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ParokiaID  uint
	Parokia    Parokia
	StartTime  sql.NullTime
	EndTime    sql.NullTime
	Details    string `gorm:"type:text"`
	LanguageID uint
	Language   Language
	WeekDay    uint
}
