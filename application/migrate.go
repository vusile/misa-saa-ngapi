package application

import (
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
)

func migrate(app *App) {

	fmt.Println("Running migrations")
	app.gorm.AutoMigrate(&model.Huduma{})
	app.gorm.AutoMigrate(&model.User{})
	app.gorm.AutoMigrate(&model.Church{})
	app.gorm.AutoMigrate(&model.Country{})
	app.gorm.AutoMigrate(&model.Jimbo{})
	app.gorm.AutoMigrate(&model.Parokia{})
	app.gorm.AutoMigrate(&model.History{})
	app.gorm.AutoMigrate(&model.Language{})
	app.gorm.AutoMigrate(&model.Timing{})
	app.gorm.AutoMigrate(&model.ModelType{})
	app.gorm.AutoMigrate(&model.WeekDay{})

	var modelTypes = []*model.ModelType{
		{Name: "huduma"},
		{Name: "jimbo"},
		{Name: "parokia"},
	}

	var weekdays = []*model.WeekDay{
		{Name: "Jumapili"},
		{Name: "Jumatatu"},
		{Name: "Jumanne"},
		{Name: "Jumatano"},
		{Name: "Alhamisi"},
		{Name: "Ijumaa"},
		{Name: "Jumamosi"},
	}

	app.gorm.Create(modelTypes)
	app.gorm.Create(weekdays)
}
