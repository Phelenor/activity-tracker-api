package activity

//go:generate go run github.com/dmarkham/enumer -type=ActivityType -json -text -output activitytype_string.go -trimprefix ActivityType -transform snake-upper
type ActivityType int

const (
	ActivityTypeUndefined ActivityType = iota
	ActivityTypeRun
	ActivityTypeWalk
	ActivityTypeCycling
	ActivityTypeOther
)
