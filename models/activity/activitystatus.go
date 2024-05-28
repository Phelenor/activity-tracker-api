package activity

//go:generate go run github.com/dmarkham/enumer -type=ActivityStatus -json -text -output activitystatus_string.go -trimprefix ActivityStatus -transform snake-upper
type ActivityStatus int

const (
	ActivityStatusUndefined ActivityStatus = iota
	ActivityStatusNotStarted
	ActivityStatusInProgress
	ActivityStatusPaused
	ActivityStatusFinished
	ActivityStatusCanceled
)
