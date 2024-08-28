package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorLanguageNotExist = errors.New("language does not exist")

type LanguageRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *LanguageRepo) Insert(ctx context.Context, language model.Language) error {

	tx := repo.Client.Create(&language)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert language %e", tx.Error)
	}

	return nil
}

func (repo *LanguageRepo) FindByID(ctx context.Context, id uint64) (model.Language, error) {

	var language model.Language

	tx := repo.Client.First(&language, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Language{}, ErrorLanguageNotExist
	}

	return language, nil
}

func (repo *LanguageRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Language{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorLanguageNotExist
	}

	return nil
}

func (repo *LanguageRepo) Update(ctx context.Context, language model.Language) error {

	tx := repo.Client.Save(language)

	if tx.Error != nil {
		return fmt.Errorf("failed to update Language %e", tx.Error)
	}

	return nil
}

type FindLanguageResult struct {
	Languages []model.Language
	Page      int
}

func (repo *LanguageRepo) FindAll(ctx context.Context, page FindAllPage) (FindLanguageResult, error) {

	var languages []model.Language

	tx := repo.Client.Scopes(Paginate(page)).Find(&languages)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindLanguageResult{}, ErrorLanguageNotExist
	}

	return FindLanguageResult{
			Languages: languages,
			Page:      page.PageNum + 1,
		},
		nil
}
