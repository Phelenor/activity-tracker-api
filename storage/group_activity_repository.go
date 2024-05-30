package storage

import (
	"activity-tracker-api/models/activity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type GroupActivityRepository interface {
	GetByIDFromDb(id string) (*activity.GroupActivity, error)
	GetByIDFromRedis(id string) (*activity.GroupActivity, error)
	GetByJoinCodeFromRedis(id string) (*activity.GroupActivity, error)
	Insert(dbActivity *activity.GroupActivity) error
	Delete(id string) error
	DeleteExpiredActivities() error
	GetByUserIdFromRedis(userId string) ([]*activity.GroupActivity, error)
}

type groupActivityRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewGroupActivityRepository(db *gorm.DB, redis *redis.Client) GroupActivityRepository {
	return &groupActivityRepo{db: db, redis: redis}
}

var ctx = context.Background()

func (repo *groupActivityRepo) GetByIDFromDb(id string) (*activity.GroupActivity, error) {
	var groupActivityDb activity.GroupActivity
	result := repo.db.First(&groupActivityDb, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &groupActivityDb, nil
}

func (repo *groupActivityRepo) GetByIDFromRedis(id string) (*activity.GroupActivity, error) {
	activityBytes, err := repo.redis.Get(ctx, id).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("no activity found with id %s", id)
	}

	var groupActivity activity.GroupActivity
	if err := json.Unmarshal(activityBytes, &groupActivity); err != nil {
		return nil, err
	}

	return &groupActivity, nil
}

func (repo *groupActivityRepo) GetByJoinCodeFromRedis(joinCode string) (*activity.GroupActivity, error) {
	idBytes, err := repo.redis.Get(ctx, joinCode).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("invalid join code: %s", joinCode)
	}

	id := string(idBytes)

	return repo.GetByIDFromRedis(id)
}

func (repo *groupActivityRepo) Insert(groupActivity *activity.GroupActivity) error {
	if groupActivity.Status == activity.ActivityStatusFinished {
		err := repo.deleteActivityRefsFromRedis(groupActivity)
		if err != nil {
			return err
		}

		return repo.db.Create(groupActivity).Error
	}

	activityJSON, err := json.Marshal(groupActivity)
	if err != nil {
		return err
	}

	err = repo.redis.Set(ctx, groupActivity.Id, activityJSON, 0).Err()
	if err != nil {
		return err
	}

	for _, userId := range groupActivity.JoinedUsers {
		err = repo.redis.SAdd(ctx, "user:"+userId+":activities", groupActivity.Id).Err()
		if err != nil {
			return err
		}
	}

	return repo.redis.Set(ctx, groupActivity.JoinCode, []byte(groupActivity.Id), 0).Err()
}

func (repo *groupActivityRepo) Delete(id string) error {
	activityBytes, err := repo.redis.Get(ctx, id).Bytes()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("no activity found with id %s", id)
	}

	var groupActivity activity.GroupActivity
	if err := json.Unmarshal(activityBytes, &groupActivity); err != nil {
		return err
	}

	return repo.deleteActivityRefsFromRedis(&groupActivity)
}

func (repo *groupActivityRepo) DeleteExpiredActivities() error {
	expiredActivities, err := repo.getExpiredActivities()
	if err != nil {
		return err
	}

	for _, groupActivity := range expiredActivities {
		err = repo.deleteActivityRefsFromRedis(groupActivity)
	}

	return err
}

func (repo *groupActivityRepo) getExpiredActivities() ([]*activity.GroupActivity, error) {
	var expiredActivities []*activity.GroupActivity
	now := time.Now().Unix()
	twoHoursAgo := now - 2*3600

	iter := repo.redis.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		id := iter.Val()
		activityJSON, err := repo.redis.Get(ctx, id).Bytes()

		if errors.Is(err, redis.Nil) {
			continue
		} else if err != nil {
			return nil, err
		}

		var groupActivity activity.GroupActivity
		if err := json.Unmarshal(activityJSON, &groupActivity); err != nil {
			continue
		}

		if groupActivity.Status == activity.ActivityStatusNotStarted && groupActivity.StartTimestamp <= twoHoursAgo {
			expiredActivities = append(expiredActivities, &groupActivity)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	log.Info("Deleting expired activities from Redis ", expiredActivities)

	return expiredActivities, nil
}

func (repo *groupActivityRepo) GetByUserIdFromRedis(userId string) ([]*activity.GroupActivity, error) {
	activityIds := repo.redis.SMembers(ctx, "user:"+userId+":activities").Val()

	var activities []*activity.GroupActivity
	for _, activityId := range activityIds {
		groupActivity, err := repo.GetByIDFromRedis(activityId)
		if err != nil {
			return nil, err
		}

		activities = append(activities, groupActivity)
	}

	return activities, nil
}

func (repo *groupActivityRepo) deleteActivityRefsFromRedis(groupActivity *activity.GroupActivity) error {
	err := repo.redis.Del(ctx, groupActivity.JoinCode).Err()
	err = repo.redis.Del(ctx, groupActivity.Id).Err()

	for _, userId := range groupActivity.JoinedUsers {
		err = repo.redis.SRem(ctx, "user:"+userId+":activities", groupActivity.Id).Err()
		if err != nil {
			return err
		}
	}

	return err
}
