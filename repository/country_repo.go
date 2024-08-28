package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorCountryNotExist = errors.New("country does not exist")

type CountryRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *CountryRepo) Insert(ctx context.Context, country model.Country) error {

	tx := repo.Client.Create(&country)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert country %e", tx.Error)
	}

	return nil
}

func (repo *CountryRepo) FindByID(ctx context.Context, id uint64) (model.Country, error) {

	var country model.Country

	tx := repo.Client.First(&country, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Country{}, ErrorCountryNotExist
	}

	return country, nil
}

func (repo *CountryRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Country{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorCountryNotExist
	}

	return nil
}

func (repo *CountryRepo) Update(ctx context.Context, country model.Country) error {

	tx := repo.Client.Save(country)

	if tx.Error != nil {
		return fmt.Errorf("failed to update country %e", tx.Error)
	}

	return nil
}

type FindCountryResult struct {
	Countries []model.Country
	Page      int
}

func (repo *CountryRepo) FindAll(ctx context.Context, page FindAllPage) (FindCountryResult, error) {

	var countries []model.Country

	tx := repo.Client.Scopes(Paginate(page)).Find(&countries)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindCountryResult{}, ErrorCountryNotExist
	}

	return FindCountryResult{
			Countries: countries,
			Page:      page.PageNum + 1,
		},
		nil
}
