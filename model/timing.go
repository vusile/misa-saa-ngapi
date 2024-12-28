package model

import (
	"time"

	"gorm.io/gorm"
)

type Timing struct {
	gorm.Model
	ID              uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ParokiaID       uint
	Parokia         Parokia
	StartTime       *time.Time `gorm:"type:TIME;null;default:null"`
	EndTime         *time.Time `gorm:"type:TIME;null;default:null"`
	Details         string     `gorm:"type:text"`
	LanguageID      uint
	Language        Language
	WeekDayID       uint
	WeekDay         WeekDay
	UserID          uint
	User            User
	HudumaID        uint
	Huduma          Huduma
	IsPublicHoliday bool
	InsertGroupBy   int64
}
