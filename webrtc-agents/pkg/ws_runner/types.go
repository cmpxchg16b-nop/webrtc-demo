package ws_runner

import (
	"context"

	pkgframing "example.com/webrtcserver/pkg/framing"
)

type WebSocketSignallingSessionRunner interface {
	Run(ctx context.Context) (chan<- pkgframing.MessagePayload, <-chan pkgframing.MessagePayload)
}
