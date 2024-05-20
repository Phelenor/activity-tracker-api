package activity

//go:generate go run github.com/dmarkham/enumer -type=HeartRateZone -json -output heartratezone_string.go -trimprefix HeartRateZone -transform snake-upper
type HeartRateZone int

const (
	HeartRateZoneUndefined HeartRateZone = iota
	HeartRateZoneAtRest
	HeartRateZoneWarmUp
	HeartRateZoneFatBurn
	HeartRateZoneAerobic
	HeartRateZoneAnaerobic
	HeartRateZoneVo2Max
)
