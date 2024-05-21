package activity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Goal struct {
	Type      ActivityGoalType        `json:"type"`
	ValueType GoalValueComparisonType `json:"valueType"`
	Label     string                  `json:"label,omitempty"`
	GoalValue float32                 `json:"value"`
}

func (goal *Goal) Value() (driver.Value, error) {
	return json.Marshal(goal)
}

func (goal *Goal) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, goal)
}

type GoalProgress struct {
	Goal         Goal    `json:"goal"`
	CurrentValue float32 `json:"currentValue"`
}

func (goalProgress *GoalProgress) Value() (driver.Value, error) {
	return json.Marshal(goalProgress)
}

func (goalProgress *GoalProgress) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, goalProgress)
}

type Goals []GoalProgress

func (g *Goals) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *Goals) Scan(value interface{}) error {
	if value == nil {
		*g = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, g)
}
