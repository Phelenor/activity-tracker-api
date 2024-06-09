package gym

type SimulatorDataSnapshot struct {
	Type          string  `json:"type"`
	Duration      int32   `json:"duration"`
	Distance      int     `json:"distance"`
	Speed         float32 `json:"speed"`
	AvgSpeed      float32 `json:"avgSpeed"`
	MaxSpeed      float32 `json:"maxSpeed"`
	HeartRate     int     `json:"heartRate"`
	AvgHeartRate  int     `json:"avgHeartRate"`
	MaxHeartRate  int     `json:"maxHeartRate"`
	ElevationGain int     `json:"elevationGain"`
	Calories      int     `json:"calories"`
}
