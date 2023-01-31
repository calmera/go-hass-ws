package go_hass_ws

import "encoding/json"

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

type ServicesCallback func(services map[string]ServiceDomain)
type ServiceDomain map[string]Service
type Service struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Fields      map[string]ServiceField `json:"fields"`
	Target      map[string]interface{}  `json:"target"`
}
type ServiceField struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Required    bool                       `json:"required"`
	Example     interface{}                `json:"example"`
	Selector    map[string]json.RawMessage `json:"selector"`
}
