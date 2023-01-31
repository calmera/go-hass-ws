package go_hass_ws

import "encoding/json"

type HassMessage struct {
	Type string `json:"type"`
}

type HassAuthenticatedMessage struct {
	HassMessage
	Id      uint64          `json:"id"`
	Event   json.RawMessage `json:"event"`
	Result  json.RawMessage `json:"result"`
	Success bool            `json:"success"`
}
