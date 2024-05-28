package storage

import (
	"activity-tracker-api/models/activity"
	"encoding/json"
	"fmt"
	"github.com/gofiber/storage/redis/v3"
	"gorm.io/gorm"
)

type GroupActivityRepository interface {
	GetByIDFromDb(id string) (*activity.GroupActivity, error)
	GetByIDFromRedis(id string) (*activity.GroupActivity, error)
	Insert(dbActivity *activity.GroupActivity) error
	Delete(id string, userId string) error
}

type groupActivityRepo struct {
	db    *gorm.DB
	redis *redis.Storage
}

func NewGroupActivityRepository(db *gorm.DB, redis *redis.Storage) GroupActivityRepository {
	return &groupActivityRepo{db: db, redis: redis}
}

func (repo *groupActivityRepo) GetByIDFromDb(id string) (*activity.GroupActivity, error) {
	var groupActivityDb activity.GroupActivity
	result := repo.db.First(&groupActivityDb, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &groupActivityDb, nil
}

func (repo *groupActivityRepo) GetByIDFromRedis(id string) (*activity.GroupActivity, error) {
	activityBytes, err := repo.redis.Get(id)

	if err != nil || activityBytes == nil {
		return nil, err
	}

	var groupActivity activity.GroupActivity
	if err := json.Unmarshal(activityBytes, &groupActivity); err != nil {
		return nil, err
	}

	return &groupActivity, nil
}

func (repo *groupActivityRepo) Insert(groupActivity *activity.GroupActivity) error {
	if groupActivity.Status == activity.ActivityStatusFinished {
		return repo.db.Create(groupActivity).Error
	}

	activityJSON, err := json.Marshal(groupActivity)
	if err != nil {
		return err
	}

	return repo.redis.Set(groupActivity.Id, activityJSON, 0)
}

func (repo *groupActivityRepo) Delete(id string, userId string) error {
	activityBytes, err := repo.redis.Get(id)

	if err != nil || activityBytes == nil {
		return err
	}

	var groupActivity activity.GroupActivity
	if err := json.Unmarshal(activityBytes, &groupActivity); err != nil {
		return err
	}

	if groupActivity.UserOwnerId == userId {
		return repo.redis.Delete(id)
	}

	return fmt.Errorf("unathorized to delete group activity %s", id)
}
