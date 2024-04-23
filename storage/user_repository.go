package storage

import (
	"activity-tracker-api/models"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id string) (*models.User, error)
	Insert(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	result := repo.db.First(&user, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (repo *userRepository) Insert(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) Update(user *models.User) error {
	return repo.db.Updates(user).Error
}

func (repo *userRepository) Delete(id string) error {
	result := repo.db.Delete(&models.User{}, "id = ?", id)
	return result.Error
}
