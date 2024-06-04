package ws

//go:generate go run github.com/dmarkham/enumer -type=ActivityControl -json -text -output activity_control_string.go -trimprefix ActivityControl -transform snake-upper
type ActivityControl int

const (
	ActivityControlUndefined ActivityControl = iota
	ActivityControlStart
	ActivityControlPause
	ActivityControlResume
	ActivityControlFinish
)
