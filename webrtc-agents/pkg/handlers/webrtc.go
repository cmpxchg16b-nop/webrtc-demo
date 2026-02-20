package handlers

import (
	"context"

	pkgframing "example.com/webrtcserver/pkg/framing"
)

type WebRTCHandler struct {
}

type SignallingServerProxy struct {
	GetIO func(ctx context.Context) (chan<- pkgframing.MessagePayload, <-chan pkgframing.MessagePayload)
}

func (handler *WebRTCHandler) Run(ctx context.Context, signallingServer SignallingServerProxy) {

}
