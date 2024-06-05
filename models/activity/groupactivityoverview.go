package activity

import "activity-tracker-api/models"

type GroupActivityOverview struct {
	Users   []models.User `json:"users"`
	OwnerId string        `json:"ownerId"`
}
