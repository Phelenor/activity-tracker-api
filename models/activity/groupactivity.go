package activity

type GroupActivity struct {
	Id             string         `json:"id" gorm:"primaryKey"`
	JoinCode       string         `json:"joinCode"`
	UserOwnerId    string         `json:"userOwnerId"`
	ActivityType   ActivityType   `json:"activityType"`
	StartTimestamp int64          `json:"startTimestamp"`
	Status         ActivityStatus `json:"status"`
	StartedUsers   []string       `json:"startedUsers" gorm:"type:jsonb"`
	ActiveUsers    []string       `json:"activeUsers" gorm:"type:jsonb"`
}
