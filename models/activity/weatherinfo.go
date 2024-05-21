package activity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type WeatherInfo struct {
	Temp        float32 `json:"temp"`
	Humidity    float32 `json:"humidity"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
}

func (weatherInfo *WeatherInfo) Value() (driver.Value, error) {
	return json.Marshal(weatherInfo)
}

func (weatherInfo *WeatherInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, weatherInfo)
}
