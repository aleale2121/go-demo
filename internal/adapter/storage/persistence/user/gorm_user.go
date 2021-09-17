package user

import (
	"template/internal/constant/errors"
	"template/internal/constant/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type userGormRepo struct {
	conn *gorm.DB
}

func InitUserPersistence(db *gorm.DB) UserStorage {
	return &userGormRepo{conn: db}
}

func (repo userGormRepo) CreateUser(usr *model.User) (*model.User, error) {
	err := repo.conn.Create(&usr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// return nil, errors
		}

		// return nil, errors.ErrorUnableToSave
	}
	return usr, nil
}

func (repo userGormRepo) DeleteUser(id uuid.UUID) error {
	err := repo.conn.Delete(&model.User{}, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// return errors.IDNotFound
		}
		// return errors.ErrorUnableToDelete
	}
	return nil
}

func (repo userGormRepo) GetUserById(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	err := repo.conn.First(user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrIDNotFound
		}
		return nil, errors.ErrUnableToFetch
	}
	return user, nil
}

func (repo userGormRepo) GetUsers() ([]model.User, error) {
	return nil, nil
}
