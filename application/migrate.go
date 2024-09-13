package application

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/vusile/misa-saa-ngapi/model"
)

func migrate(app *App) {

	fmt.Println("Running migrations")
	app.gorm.AutoMigrate(&model.Huduma{})
	app.gorm.AutoMigrate(&model.User{})
	app.gorm.AutoMigrate(&model.Church{})
	app.gorm.AutoMigrate(&model.Jimbo{})
	app.gorm.AutoMigrate(&model.Parokia{})
	app.gorm.AutoMigrate(&model.History{})
	app.gorm.AutoMigrate(&model.Language{})
	app.gorm.AutoMigrate(&model.Timing{})
	app.gorm.AutoMigrate(&model.ModelType{})

	var modelTypes = []*model.ModelType{
		{Name: "huduma"},
		{Name: "jimbo"},
		{Name: "parokia"},
	}

	app.gorm.Create(modelTypes)

	var user model.User
	middleName := "Terence"
	birthday := time.Now()
	app.gorm.FirstOrCreate(&user, model.User{
		FirstName:    "Vusile",
		MiddleName:   &middleName,
		LastName:     "Silonda",
		Email:        new(string),
		Birthday:     &birthday,
		MemberNumber: sql.NullString{},
		ActivatedAt:  sql.NullTime{},
		Password:     "",
	})
}
