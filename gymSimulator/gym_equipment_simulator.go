package gymSimulator

import (
	"activity-tracker-api/models/gym"
	"math/rand"
	"sync"
)

type GymEquipmentSimulator struct {
	mu         sync.Mutex
	active     bool
	finished   bool
	data       gym.SimulatorDataSnapshot
	totalSpeed float32
	totalHR    int
	iterations int
}

func NewGymEquipmentSimulator() *GymEquipmentSimulator {
	return &GymEquipmentSimulator{
		data: gym.SimulatorDataSnapshot{Type: "gym_data_snapshot"},
	}
}

func (s *GymEquipmentSimulator) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = true
}

func (s *GymEquipmentSimulator) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
}

func (s *GymEquipmentSimulator) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = true
}

func (s *GymEquipmentSimulator) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active = false
	s.finished = true
}

func (s *GymEquipmentSimulator) IsActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active
}

func (s *GymEquipmentSimulator) IsFinished() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.finished
}

func (s *GymEquipmentSimulator) GenerateDataSnapshot() gym.SimulatorDataSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.Duration += 1

	s.data.Distance += int(s.data.Speed * 1000 / 3600)

	s.data.HeartRate = 120 + rand.Intn(10) - 5

	if s.data.HeartRate > s.data.MaxHeartRate {
		s.data.MaxHeartRate = s.data.HeartRate
	}

	s.data.Speed = 8.0 + rand.Float32()*2 - 1

	if s.data.Speed > s.data.MaxSpeed {
		s.data.MaxSpeed = s.data.Speed
	}

	if s.data.Duration%30 == 0 {
		s.data.ElevationGain += rand.Intn(2)
	}

	if s.data.Duration%6 == 0 {
		s.data.Calories += rand.Intn(3)
	}

	s.totalSpeed += s.data.Speed
	s.totalHR += s.data.HeartRate
	s.iterations++

	s.data.AvgSpeed = s.totalSpeed / float32(s.iterations)
	s.data.AvgHeartRate = s.totalHR / s.iterations

	return s.data
}
