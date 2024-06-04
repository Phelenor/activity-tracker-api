package ws

import (
	"activity-tracker-api/models/activity"
)

type UserDataUpdate struct {
	UserId          string  `json:"userId"`
	UserDisplayName string  `json:"userDisplayName"`
	UserImageUrl    string  `json:"userImageUrl,omitempty"`
	Type            string  `json:"type"`
	Lat             float32 `json:"lat"`
	Long            float32 `json:"long"`
	Distance        int     `json:"distance"`
	HeartRate       int     `json:"heartRate"`
	Speed           float32 `json:"speed"`
}

type ControlAction struct {
	Action ActivityControl `json:"action"`
}

type UserFinish struct {
	UserId          string                 `json:"userId"`
	DurationSeconds int32                  `json:"durationSeconds"`
	Activity        activity.GroupActivity `json:"activity"`
}

type ActivityUpdate struct {
	Type     string                 `json:"type"`
	Activity activity.GroupActivity `json:"activity"`
}
