package models

type Activity struct {
	ActivityType              ActivityType              `json:"activityType"`
	StartTimestamp            int64                     `json:"startTimestamp"`
	EndTimestamp              int64                     `json:"endTimestamp"`
	DistanceMeters            int                       `json:"distanceMeters"`
	DurationSeconds           int64                     `json:"durationSeconds"`
	AvgSpeedKmh               float32                   `json:"avgSpeedKmh"`
	AvgHeartRate              int                       `json:"avgHeartRate"`
	Calories                  int                       `json:"calories"`
	Elevation                 int                       `json:"elevation"`
	Weather                   *ActivityWeatherInfo      `json:"weather"`
	HeartRateZoneDistribution map[HeartRateZone]float32 `json:"heartRateZoneDistribution"`
	Goals                     []ActivityGoalProgress    `json:"goals"`
}

type ActivityWeatherInfo struct {
	Temp        float32 `json:"temp"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
}

type ActivityType string

const (
	RUN     ActivityType = "RUN"
	WALK                 = "WALK"
	CYCLING              = "CYCLING"
	OTHER                = "OTHER"
)

type HeartRateZone string

const (
	AT_REST   HeartRateZone = "AT_REST"
	WARM_UP                 = "WARM_UP"
	FAT_BURN                = "FAT_BURN"
	AEROBIC                 = "AEROBIC"
	ANAEROBIC               = "ANAEROBIC"
	VO2_MAX                 = "VO2_MAX"
)

type GoalValueComparisonType string

const (
	GREATER GoalValueComparisonType = "GREATER"
	LESS                            = "LESS"
)

type ActivityGoalType string

const (
	DISTANCE       ActivityGoalType = "DISTANCE"
	DURATION                        = "DURATION"
	CALORIES                        = "CALORIES"
	AVG_HEART_RATE                  = "AVG_HEART_RATE"
	AVG_SPEED                       = "AVG_SPEED"
	AVG_PACE                        = "AVG_PACE"
	IN_HR_ZONE                      = "IN_HR_ZONE"
	BELOW_HR_ZONE                   = "BELOW_HR_ZONE"
	ABOVE_HR_ZONE                   = "ABOVE_HR_ZONE"
)

type ActivityGoal struct {
	Type      ActivityGoalType        `json:"type"`
	ValueType GoalValueComparisonType `json:"valueType"`
	Label     string                  `json:"label,omitempty"`
	Value     float32                 `json:"value"`
}

type ActivityGoalProgress struct {
	Goal         ActivityGoal `json:"goal"`
	CurrentValue float32      `json:"currentValue"`
}
