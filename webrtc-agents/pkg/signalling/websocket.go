package signalling

import (
	"context"

	pkgframing "example.com/webrtcserver/pkg/framing"
)

// WebSocketProxy implements SignallingServerProxy using WebSocket
type WebSocketProxy struct {
	DataChan chan pkgframing.MessagePayload
	TxChan   chan pkgframing.MessagePayload
}

// Send sends a message to the signalling server
func (proxy *WebSocketProxy) Send(ctx context.Context, msg pkgframing.MessagePayload) error {
	proxy.TxChan <- msg
	return nil
}

// Receive returns a channel for receiving messages from the signalling server
func (proxy *WebSocketProxy) Receive() <-chan pkgframing.MessagePayload {
	return proxy.DataChan
}
