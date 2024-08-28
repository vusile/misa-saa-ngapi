package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorParokiaNotExist = errors.New("parokia does not exist")

type ParokiaRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *ParokiaRepo) Insert(ctx context.Context, parokia model.Parokia) error {

	tx := repo.Client.Create(&parokia)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert parokia %e", tx.Error)
	}

	return nil
}

func (repo *ParokiaRepo) FindByID(ctx context.Context, id uint64) (model.Parokia, error) {

	var parokia model.Parokia

	tx := repo.Client.Preload("Jimbo").First(&parokia, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Parokia{}, ErrorParokiaNotExist
	}

	return parokia, nil
}

func (repo *ParokiaRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Parokia{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorNotExist
	}

	return nil
}

func (repo *ParokiaRepo) Update(ctx context.Context, parokia model.Parokia) error {

	tx := repo.Client.Save(parokia)

	if tx.Error != nil {
		return fmt.Errorf("failed to update Parokia %e", tx.Error)
	}

	return nil
}

type FindParokiaResult struct {
	Parokia []model.Parokia
	Page    int
}

func (repo *ParokiaRepo) FindAll(ctx context.Context, page FindAllPage) (FindParokiaResult, error) {

	var parokia []model.Parokia

	tx := repo.Client.Scopes(Paginate(page)).Preload("Jimbo").Find(&parokia)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindParokiaResult{}, ErrorNotExist
	}

	return FindParokiaResult{
			Parokia: parokia,
			Page:    page.PageNum + 1,
		},
		nil
}
