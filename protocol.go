package go_hass_ws

type HassMessage struct {
	Type string `json:"type"`
}

type HassAuthenticatedMessage struct {
	HassMessage
	Id string `json:"id"`
}
