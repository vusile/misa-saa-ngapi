package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorTimingNotExist = errors.New("timing does not exist")

type TimingRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *TimingRepo) Insert(ctx context.Context, timing []model.Timing) error {

	tx := repo.Client.Create(&timing)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert timing %e", tx.Error)
	}

	return nil
}

func (repo *TimingRepo) FindByID(ctx context.Context, id uint64) (model.Timing, error) {

	var timing model.Timing

	tx := repo.Client.First(&timing, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.Timing{}, ErrorTimingNotExist
	}

	return timing, nil
}

func (repo *TimingRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.Timing{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorTimingNotExist
	}

	return nil
}

func (repo *TimingRepo) Update(ctx context.Context, timing model.Timing) error {

	tx := repo.Client.Save(timing)

	if tx.Error != nil {
		return fmt.Errorf("failed to update timing %e", tx.Error)
	}

	return nil
}

type FindTimingResult struct {
	Timings []model.Timing
	Page    int
}

func (repo *TimingRepo) FindAll(ctx context.Context, page FindAllPage) (FindTimingResult, error) {

	var timings []model.Timing

	tx := repo.Client.Scopes(Paginate(page)).Find(&timings)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindTimingResult{}, ErrorTimingNotExist
	}

	return FindTimingResult{
			Timings: timings,
			Page:    page.PageNum + 1,
		},
		nil
}

func (repo *TimingRepo) FindByParishId(ctx context.Context, parokiaID uint64, page FindAllPage) (FindTimingResult, error) {

	var timings []model.Timing

	tx := repo.Client.InnerJoins("Huduma").InnerJoins("Language").InnerJoins("WeekDay").Find(&timings, "parokia_id = ?", parokiaID)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindTimingResult{}, ErrorTimingNotExist
	}

	return FindTimingResult{
			Timings: timings,
			Page:    page.PageNum + 1,
		},
		nil
}
