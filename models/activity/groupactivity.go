package activity

type GroupActivity struct {
	Id             string         `json:"id" gorm:"primaryKey"`
	JoinCode       string         `json:"joinCode"`
	UserOwnerId    string         `json:"userOwnerId"`
	UserOwnerName  string         `json:"userOwnerName"`
	ActivityType   ActivityType   `json:"activityType"`
	StartTimestamp int64          `json:"startTimestamp"`
	Status         ActivityStatus `json:"status"`
	JoinedUsers    []string       `json:"joinedUsers" gorm:"type:jsonb"`
	StartedUsers   []string       `json:"startedUsers" gorm:"type:jsonb"`
	ActiveUsers    []string       `json:"activeUsers" gorm:"type:jsonb"`
}
