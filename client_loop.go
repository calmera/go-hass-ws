package go_hass_ws

import (
	"context"
	"nhooyr.io/websocket/wsjson"
)

func (hc *HassClient) loop(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			var msg HassAuthenticatedMessage
			if err := wsjson.Read(context.Background(), hc.c, &msg); err != nil {
				if hc.Errors != nil {
					hc.Errors <- err
				}

				continue
			}

			cb, fnd := hc.callbacks[msg.Id]
			if !fnd {
				continue
			}

			shouldRemove, err := cb(msg)
			if err != nil {
				if hc.Errors != nil {
					hc.Errors <- err
				}

				continue
			}

			if shouldRemove {
				hc.wg.Add(-1)
				delete(hc.callbacks, msg.Id)
			}
		}
	}
}

func (hc *HassClient) WaitUntilAllHandled() {
	hc.wg.Wait()
}
