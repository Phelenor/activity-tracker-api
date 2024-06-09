package storage

import (
	"activity-tracker-api/models/gym"
	"gorm.io/gorm"
)

type GymAccountRepository interface {
	GetByID(id string) (*gym.GymAccount, error)
	GetByEmail(email string) (*gym.GymAccount, error)
	Insert(GymAccount *gym.GymAccount) error
	Update(GymAccount *gym.GymAccount) error
	Delete(id string) error
}

type gymAccountRepository struct {
	db *gorm.DB
}

func NewGymAccountRepository(db *gorm.DB) GymAccountRepository {
	return &gymAccountRepository{db: db}
}

func (repo *gymAccountRepository) GetByID(id string) (*gym.GymAccount, error) {
	var GymAccount gym.GymAccount
	result := repo.db.First(&GymAccount, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &GymAccount, nil
}

func (repo *gymAccountRepository) GetByEmail(email string) (*gym.GymAccount, error) {
	var GymAccount gym.GymAccount
	result := repo.db.First(&GymAccount, "email = ?", email)

	if result.Error != nil {
		return nil, result.Error
	}

	return &GymAccount, nil
}

func (repo *gymAccountRepository) Insert(GymAccount *gym.GymAccount) error {
	return repo.db.Create(GymAccount).Error
}

func (repo *gymAccountRepository) Update(GymAccount *gym.GymAccount) error {
	return repo.db.Updates(GymAccount).Error
}

func (repo *gymAccountRepository) Delete(id string) error {
	result := repo.db.Delete(&gym.GymAccount{}, "id = ?", id)
	return result.Error
}
