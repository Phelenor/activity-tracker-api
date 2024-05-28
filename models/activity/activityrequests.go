package activity

type CreateGroupActivityRequest struct {
	ActivityType   ActivityType `json:"activityType"`
	StartTimestamp int64        `json:"startTimestamp"`
}
