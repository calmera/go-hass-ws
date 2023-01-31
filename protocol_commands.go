package go_hass_ws

type CommandMessage struct {
	Id   uint64 `json:"id"`
	Type string `json:"type"`
}

type StatesCallback func(states map[string]State)

type State struct {
	EntityId    string                 `json:"entity_id"`
	LastChanged string                 `json:"last_changed"`
	Value       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
}

type ConfigCallback func(config HassConfig)
type HassConfig map[string]interface{}

type ServicesCallback func(services map[string]Service)
type Service struct {
	Domain   string   `json:"domain"`
	Services []string `json:"services"`
}
