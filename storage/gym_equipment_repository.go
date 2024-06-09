package storage

import (
	"activity-tracker-api/models/gym"
	"gorm.io/gorm"
)

type GymEquipmentRepository interface {
	GetForUserId(id string) ([]gym.Equipment, error)
	Insert(equipment *gym.Equipment) error
	Update(equipment *gym.Equipment) error
	Delete(id string) error
}

type gymEquipmentRepository struct {
	db *gorm.DB
}

func NewGymEquipmentRepository(db *gorm.DB) GymEquipmentRepository {
	return &gymEquipmentRepository{db: db}
}

func (repo *gymEquipmentRepository) GetForUserId(id string) ([]gym.Equipment, error) {
	var equipment []gym.Equipment
	result := repo.db.Find(&equipment, "owner_id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return equipment, nil
}

func (repo *gymEquipmentRepository) Insert(equipment *gym.Equipment) error {
	return repo.db.Create(equipment).Error
}

func (repo *gymEquipmentRepository) Update(equipment *gym.Equipment) error {
	return repo.db.Updates(equipment).Error
}

func (repo *gymEquipmentRepository) Delete(id string) error {
	result := repo.db.Delete(&gym.Equipment{}, "id = ?", id)
	return result.Error
}
