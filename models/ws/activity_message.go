package ws

import "encoding/json"

type ActivityMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ConnectMessage struct {
	Data string `json:"data"`
}

type DataUpdate struct {
	Data string `json:"data"`
}

type StatusChange struct {
	Data string `json:"data"`
}
