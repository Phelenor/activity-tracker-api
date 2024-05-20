package activity

//go:generate go run github.com/dmarkham/enumer -type=GoalValueComparisonType -json --output goalvaluecomparisontype_string.go -trimprefix GoalValueComparisonType -transform snake-upper
type GoalValueComparisonType int

const (
	GoalValueComparisonTypeUndefined GoalValueComparisonType = iota
	GoalValueComparisonTypeGreater
	GoalValueComparisonTypeLess
)
