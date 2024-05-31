package activity

type Activity struct {
	Id                        string                    `json:"id"`
	GroupActivityId           string                    `json:"groupActivityId,omitempty"`
	ActivityType              ActivityType              `json:"activityType"`
	StartTimestamp            int64                     `json:"startTimestamp"`
	EndTimestamp              int64                     `json:"endTimestamp"`
	DistanceMeters            int                       `json:"distanceMeters"`
	DurationSeconds           int64                     `json:"durationSeconds"`
	AvgSpeedKmh               float32                   `json:"avgSpeedKmh"`
	MaxSpeedKmh               float32                   `json:"maxSpeedKmh"`
	AvgHeartRate              int                       `json:"avgHeartRate"`
	MaxHeartRate              int                       `json:"maxHeartRate"`
	Calories                  int                       `json:"calories"`
	Elevation                 int                       `json:"elevation"`
	Weather                   *WeatherInfo              `json:"weather"`
	HeartRateZoneDistribution HeartRateZoneDistribution `json:"heartRateZoneDistribution"`
	Goals                     Goals                     `json:"goals"`
	ImageUrl                  string                    `json:"imageUrl"`
}

type DbActivity struct {
	Id                        string                    `json:"id" gorm:"primaryKey"`
	GroupActivityId           string                    `json:"groupActivityId,omitempty"`
	UserId                    string                    `json:"userId"`
	ActivityType              ActivityType              `json:"activityType"`
	StartTimestamp            int64                     `json:"startTimestamp"`
	EndTimestamp              int64                     `json:"endTimestamp"`
	DistanceMeters            int                       `json:"distanceMeters"`
	DurationSeconds           int64                     `json:"durationSeconds"`
	AvgSpeedKmh               float32                   `json:"avgSpeedKmh"`
	MaxSpeedKmh               float32                   `json:"maxSpeedKmh"`
	AvgHeartRate              int                       `json:"avgHeartRate"`
	MaxHeartRate              int                       `json:"maxHeartRate"`
	Calories                  int                       `json:"calories"`
	Elevation                 int                       `json:"elevation"`
	Weather                   *WeatherInfo              `json:"weather" gorm:"type:jsonb"`
	HeartRateZoneDistribution HeartRateZoneDistribution `json:"heartRateZoneDistribution" gorm:"type:jsonb"`
	Goals                     Goals                     `json:"goals" gorm:"type:jsonb"`
}
