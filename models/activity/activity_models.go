package activity

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
