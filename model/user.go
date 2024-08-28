package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	FirstName    string
	MiddleName   *string
	LastName     string
	Email        *string
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	Password     string
}
