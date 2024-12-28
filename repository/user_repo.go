package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vusile/misa-saa-ngapi/model"
	"gorm.io/gorm"
)

var ErrorUserNotExist = errors.New("user does not exist")

type UserRepo struct {
	Client *gorm.DB
	DB     *sql.DB
}

func (repo *UserRepo) Insert(ctx context.Context, user model.User) (uint, error) {

	tx := repo.Client.Create(&user)

	if tx.Error != nil {
		return 0, fmt.Errorf("failed to insert user %e", tx.Error)
	}

	return user.ID, nil
}

func (repo *UserRepo) FindByID(ctx context.Context, id uint64) (model.User, error) {

	var user model.User

	tx := repo.Client.First(&user, "id = ?", id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return model.User{}, ErrorUserNotExist
	}

	return user, nil
}

func (repo *UserRepo) DeleteByID(ctx context.Context, id uint64) error {

	tx := repo.Client.Delete(&model.User{}, id)

	if tx.Error != nil || tx.RowsAffected == 0 {
		return ErrorUserNotExist
	}

	return nil
}

func (repo *UserRepo) Update(ctx context.Context, user model.User) error {

	tx := repo.Client.Save(user)

	if tx.Error != nil {
		return fmt.Errorf("failed to update user %e", tx.Error)
	}

	return nil
}

type FindUserResult struct {
	Users []model.User
	Page  int
}

func (repo *UserRepo) FindAll(ctx context.Context, page FindAllPage) (FindUserResult, error) {

	var users []model.User

	tx := repo.Client.Scopes(Paginate(page)).Find(&users)

	if tx.Error != nil {
		fmt.Println("an error occured while querying", tx.Error)
		return FindUserResult{}, ErrorUserNotExist
	}

	return FindUserResult{
			Users: users,
			Page:  page.PageNum + 1,
		},
		nil
}
