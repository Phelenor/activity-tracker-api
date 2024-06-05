package activity

import "activity-tracker-api/models/util"

type GroupActivity struct {
	Id             string         `json:"id" gorm:"primaryKey"`
	JoinCode       string         `json:"joinCode"`
	UserOwnerId    string         `json:"userOwnerId"`
	UserOwnerName  string         `json:"userOwnerName"`
	ActivityType   ActivityType   `json:"activityType"`
	StartTimestamp int64          `json:"startTimestamp"`
	Status         ActivityStatus `json:"status"`
	JoinedUsers    []string       `json:"joinedUsers" gorm:"type:jsonb"`
	ConnectedUsers []string       `json:"connectedUsers" gorm:"type:jsonb"`
	ActiveUsers    []string       `json:"activeUsers" gorm:"type:jsonb"`
	FinishedUsers  []string       `json:"finishedUsers" gorm:"type:jsonb"`
}

type DbGroupActivity struct {
	Id             string           `json:"id" gorm:"primaryKey"`
	JoinCode       string           `json:"joinCode"`
	UserOwnerId    string           `json:"userOwnerId"`
	UserOwnerName  string           `json:"userOwnerName"`
	ActivityType   ActivityType     `json:"activityType"`
	StartTimestamp int64            `json:"startTimestamp"`
	Status         ActivityStatus   `json:"status"`
	JoinedUsers    util.StringSlice `json:"joinedUsers" gorm:"type:jsonb"`
	ConnectedUsers util.StringSlice `json:"connectedUsers" gorm:"type:jsonb"`
	ActiveUsers    util.StringSlice `json:"activeUsers" gorm:"type:jsonb"`
	FinishedUsers  util.StringSlice `json:"finishedUsers" gorm:"type:jsonb"`
}
