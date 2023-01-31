package go_hass_ws

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
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

	return &HassClient{
		c:                 c,
		nextInteractionId: 0,
		Version:           authRequired.Version,
	}, nil
}

type HassClient struct {
	c                 *websocket.Conn
	nextInteractionId uint64
	Version           string
}

func (hc *HassClient) Close() error {
	return hc.c.Close(websocket.StatusNormalClosure, "")
}
