package storage

import (
	"activity-tracker-api/models"
	"gorm.io/gorm"
)

type GymRepository interface {
	GetByID(id string) (*models.GymAccount, error)
	GetByEmail(email string) (*models.GymAccount, error)
	Insert(GymAccount *models.GymAccount) error
	Update(GymAccount *models.GymAccount) error
	Delete(id string) error
}

type gymRepository struct {
	db *gorm.DB
}

func NewGymRepository(db *gorm.DB) GymRepository {
	return &gymRepository{db: db}
}

func (repo *gymRepository) GetByID(id string) (*models.GymAccount, error) {
	var GymAccount models.GymAccount
	result := repo.db.First(&GymAccount, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &GymAccount, nil
}

func (repo *gymRepository) GetByEmail(email string) (*models.GymAccount, error) {
	var GymAccount models.GymAccount
	result := repo.db.First(&GymAccount, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &GymAccount, nil
}

func (repo *gymRepository) Insert(GymAccount *models.GymAccount) error {
	return repo.db.Create(GymAccount).Error
}

func (repo *gymRepository) Update(GymAccount *models.GymAccount) error {
	return repo.db.Updates(GymAccount).Error
}

func (repo *gymRepository) Delete(id string) error {
	result := repo.db.Delete(&models.GymAccount{}, "id = ?", id)
	return result.Error
}
