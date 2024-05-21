package activity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type HeartRateZoneDistribution map[string]float32

func (h HeartRateZoneDistribution) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HeartRateZoneDistribution) Scan(value interface{}) error {
	if value == nil {
		*h = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, h)
}
