package model

type Status struct {
	ID    uint `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title string
	Model string
}
