package activity

//go:generate go run github.com/dmarkham/enumer -type=ActivityGoalType -json -output activitygoaltype_string.go -trimprefix ActivityGoalType -transform snake-upper
type ActivityGoalType int

const (
	ActivityGoalTypeUndefined ActivityGoalType = iota
	ActivityGoalTypeDistance
	ActivityGoalTypeDuration
	ActivityGoalTypeCalories
	ActivityGoalTypeAvgHeartRate
	ActivityGoalTypeAvgSpeed
	ActivityGoalTypeAvgPace
	ActivityGoalTypeInHrZone
	ActivityGoalTypeBelowHrZone
	ActivityGoalTypeAboveHrZone
)
