package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorChurchNotExist = errors.New("church does not exist")

type ChurchRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *ChurchRepo) Insert(ctx context.Context, church model.Church) error {

	tx := repo.Client.Create(&church)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert church %e", tx.Error)
	}

	return nil
}

func (repo *ChurchRepo) FindByID(ctx context.Context, id uint64) (model.Church, error) {

	var church model.Church

	tx := repo.Client.First(&church, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Church{}, ErrorChurchNotExist
	}

	return church, nil
}

func (repo *ChurchRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Church{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorChurchNotExist
	}

	return nil
}

func (repo *ChurchRepo) Update(ctx context.Context, church model.Church) error {

	tx := repo.Client.Save(church)

	if tx.Error != nil {
		return fmt.Errorf("failed to update church %e", tx.Error)
	}

	return nil
}

type FindChurchResult struct {
	Churches []model.Church
	Page     int
}

func (repo *ChurchRepo) FindAll(ctx context.Context, page FindAllPage) (FindChurchResult, error) {

	var churches []model.Church

	tx := repo.Client.Scopes(Paginate(page)).Find(&churches)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindChurchResult{}, ErrorChurchNotExist
	}

	return FindChurchResult{
			Churches: churches,
			Page:     page.PageNum + 1,
		},
		nil
}
