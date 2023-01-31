package main

import (
	"fmt"
	go_hass_ws "github.com/calmera/go-hass-ws"
)

var conf = go_hass_ws.Config{
	Url:   "ws://localhost:8123/api/websocket",
	Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIwOWQ2NDYxMzNhMzU0NmMwOWExNDYxOTEzYmUzMDU2ZCIsImlhdCI6MTY3NTE4MTg1NCwiZXhwIjoxOTkwNTQxODU0fQ.SSTRnVklTVgmH82-ndv1M01BpnCNZf23-m_eORL84Oo",
}

func main() {
	hass, err := go_hass_ws.Connect(conf)
	if err != nil {
		panic(err)
	}
	defer hass.Close()

	fmt.Printf("Connected to HASS version %s!", hass.Version)
}
