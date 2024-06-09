package storage

import (
	"activity-tracker-api/models/gym"
	"gorm.io/gorm"
)

type GymEquipmentRepository interface {
	GetForUserId(id string, name string) ([]gym.GymEquipment, error)
	GetById(id string) (*gym.GymEquipment, error)
	Insert(equipment *gym.GymEquipment) error
	Update(equipment *gym.GymEquipment) error
	Delete(id string) error
}

type gymEquipmentRepository struct {
	db *gorm.DB
}

func NewGymEquipmentRepository(db *gorm.DB) GymEquipmentRepository {
	return &gymEquipmentRepository{db: db}
}

func (repo *gymEquipmentRepository) GetForUserId(id, name string) ([]gym.GymEquipment, error) {
	var equipment []gym.GymEquipment

	query := repo.db.Where("owner_id = ?", id)
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	result := query.Find(&equipment)
	if result.Error != nil {
		return nil, result.Error
	}

	return equipment, nil
}

func (repo *gymEquipmentRepository) GetById(id string) (*gym.GymEquipment, error) {
	var equipment gym.GymEquipment
	result := repo.db.First(&equipment, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &equipment, nil
}

func (repo *gymEquipmentRepository) Insert(equipment *gym.GymEquipment) error {
	return repo.db.Create(equipment).Error
}

func (repo *gymEquipmentRepository) Update(equipment *gym.GymEquipment) error {
	return repo.db.Updates(equipment).Error
}

func (repo *gymEquipmentRepository) Delete(id string) error {
	result := repo.db.Delete(&gym.GymEquipment{}, "id = ?", id)
	return result.Error
}
