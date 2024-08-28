package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorNotExist = errors.New("jimbo does not exist")

type JimboRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *JimboRepo) Insert(ctx context.Context, jimbo model.Jimbo) error {

	tx := repo.Client.Create(&jimbo)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert jimbo %e", tx.Error)
	}

	return nil
}

func (repo *JimboRepo) FindByID(ctx context.Context, id uint64) (model.Jimbo, error) {

	var jimbo model.Jimbo

	tx := repo.Client.First(&jimbo, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Jimbo{}, ErrorNotExist
	}

	return jimbo, nil
}

func (repo *JimboRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Jimbo{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorNotExist
	}

	return nil
}

func (repo *JimboRepo) Update(ctx context.Context, jimbo model.Jimbo) error {

	tx := repo.Client.Save(jimbo)

	if tx.Error != nil {
		return fmt.Errorf("failed to update jimbo %e", tx.Error)
	}

	return nil
}

type FindResult struct {
	Majimbo []model.Jimbo
	Page    int
}

func (repo *JimboRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {

	var majimbo []model.Jimbo

	tx := repo.Client.Scopes(Paginate(page)).Find(&majimbo)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindResult{}, ErrorNotExist
	}

	return FindResult{
			Majimbo: majimbo,
			Page:    page.PageNum + 1,
		},
		nil
}
