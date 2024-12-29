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

	var countries = []*model.Country{
		{Name: "Tanzania", CountryCode: "+255"},
		{Name: "Kenya", CountryCode: "+254"},
	}

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

	var church = []*model.Church{
		{Name: "Catholic"},
	}

	var languages = []*model.Language{
		{Name: "Kiswahili"},
		{Name: "Kiingereza"},
		{Name: "Kichaga"},
		{Name: "Kihaya"},
		{Name: "Kisukuma"},
	}

	app.gorm.Create(modelTypes)
	app.gorm.Create(weekdays)
	app.gorm.Create(countries)
	app.gorm.Create(church)
	app.gorm.Create(languages)
}
