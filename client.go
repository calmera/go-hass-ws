package go_hass_ws

import (
	"context"
	"encoding/json"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sync"
	"sync/atomic"
)

type Config struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func Connect(conf Config) (*HassClient, error) {
	ctx := context.Background()

	c, _, err := websocket.Dial(ctx, conf.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", conf.Url, err)
	}

	// -- we expect the server to send us an `auth_required` message
	authRequired, err := ExpectAuthRequired(ctx, c)
	if err != nil {
		defer c.Close(websocket.StatusInternalError, err.Error())
		return nil, err
	}

	// -- we will send the server the token back as a response
	if err := Authenticate(ctx, c, conf.Token); err != nil {
		defer c.Close(websocket.StatusInternalError, err.Error())
		return nil, err
	}

	hass := &HassClient{
		c:                 c,
		nextInteractionId: 0,
		Version:           authRequired.Version,
		callbacks:         map[uint64]Callback{},
		done:              make(chan struct{}),
	}

	go hass.loop(hass.done)

	return hass, nil
}

type Callback func(msg HassAuthenticatedMessage) (bool, error)

type HassClient struct {
	c                 *websocket.Conn
	nextInteractionId uint64
	Version           string
	Errors            chan error
	callbacks         map[uint64]Callback
	wg                sync.WaitGroup
	done              chan struct{}
}

func (hc *HassClient) Close() error {
	close(hc.done)
	return hc.c.Close(websocket.StatusNormalClosure, "")
}

func (hc *HassClient) call_command(kind string, cb Callback) error {
	// -- get the id for our request
	id := atomic.AddUint64(&hc.nextInteractionId, 1)

	// -- construct the message
	cmd := CommandMessage{
		Id:   id,
		Type: kind,
	}

	// -- register a callback
	hc.callbacks[id] = cb
	hc.wg.Add(1)

	return wsjson.Write(context.Background(), hc.c, cmd)
}

func (hc *HassClient) GetStates(cb StatesCallback) error {
	return hc.call_command("get_states", func(msg HassAuthenticatedMessage) (bool, error) {
		if !msg.Success {
			return true, fmt.Errorf("request failed")
		}

		var states []State
		if err := json.Unmarshal(msg.Result, &states); err != nil {
			return true, err
		}

		results := map[string]State{}
		for _, state := range states {
			results[state.EntityId] = state
		}

		cb(results)

		return true, nil
	})
}

func (hc *HassClient) GetConfig(cb ConfigCallback) error {
	return hc.call_command("get_config", func(msg HassAuthenticatedMessage) (bool, error) {
		if !msg.Success {
			return true, fmt.Errorf("request failed")
		}

		var config HassConfig
		if err := json.Unmarshal(msg.Result, &config); err != nil {
			return true, err
		}

		cb(config)

		return true, nil
	})
}

func (hc *HassClient) GetServices(cb ServicesCallback) error {
	return hc.call_command("get_services", func(msg HassAuthenticatedMessage) (bool, error) {
		if !msg.Success {
			return true, fmt.Errorf("request failed")
		}

		var services map[string]ServiceDomain
		if err := json.Unmarshal(msg.Result, &services); err != nil {
			return true, err
		}

		cb(services)

		return true, nil
	})
}
