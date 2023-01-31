package go_hass_ws

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type AuthRequiredMessage struct {
	HassMessage
	Version string `json:"ha_version"`
}

func ExpectAuthRequired(ctx context.Context, c *websocket.Conn) (*AuthRequiredMessage, error) {
	var result AuthRequiredMessage

	if err := wsjson.Read(ctx, c, &result); err != nil {
		return nil, fmt.Errorf("unable to read 'authRequired' message: %w", err)
	}

	return &result, nil
}

type AuthMessage struct {
	HassMessage
	Token string `json:"access_token"`
}

type AuthResponseMessage struct {
	HassMessage
	Version string `json:"ha_version"` // only filled in if 'auth_ok'
	Message string `json:"message"`    // only filled in if 'auth_invalid'
}

func Authenticate(ctx context.Context, c *websocket.Conn, token string) error {
	msg := AuthMessage{
		HassMessage: HassMessage{
			Type: "auth",
		},
		Token: token,
	}

	if err := wsjson.Write(ctx, c, msg); err != nil {
		return fmt.Errorf("unable to send authorization request: %w", err)
	}

	var response AuthResponseMessage
	if err := wsjson.Read(ctx, c, &response); err != nil {
		return fmt.Errorf("unable to read authorization response: %w", err)
	}

	if response.Type == "auth_invalid" {
		return fmt.Errorf(response.Message)
	}

	if response.Type == "auth_ok" {
		return nil
	}

	return fmt.Errorf("protocol error: received a %s message while expecting an 'auth_ok' or 'auth_invalid'", response.Type)
}
