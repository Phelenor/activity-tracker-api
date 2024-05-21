package storage

import (
	"activity-tracker-api/models/activity"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	GetByID(id string) (*activity.DbActivity, error)
	GetForUserId(id string) ([]activity.DbActivity, error)
	Insert(dbActivity *activity.DbActivity) error
	Update(dbActivity *activity.DbActivity) error
	Delete(id string, userId string) error
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (repo *activityRepository) GetByID(id string) (*activity.DbActivity, error) {
	var dbActivity activity.DbActivity
	result := repo.db.First(&dbActivity, "userId = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &dbActivity, nil
}

func (repo *activityRepository) GetForUserId(id string) ([]activity.DbActivity, error) {
	var dbActivities []activity.DbActivity
	result := repo.db.Find(&dbActivities, "user_id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return dbActivities, nil
}

func (repo *activityRepository) Insert(activity *activity.DbActivity) error {
	return repo.db.Create(activity).Error
}

func (repo *activityRepository) Update(activity *activity.DbActivity) error {
	return repo.db.Updates(activity).Error
}

func (repo *activityRepository) Delete(id string, userId string) error {
	result := repo.db.Where("user_id = ?", userId).Where("id = ?", id).Delete(&activity.DbActivity{})
	return result.Error
}
