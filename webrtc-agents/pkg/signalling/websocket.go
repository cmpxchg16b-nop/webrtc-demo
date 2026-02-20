package signalling

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	pkgframing "example.com/webrtcserver/pkg/framing"

	"github.com/gorilla/websocket"

	"time"

	pkgproxy "webrtc-agents/pkg/proxy"

	pkgconnreg "example.com/webrtcserver/pkg/connreg"
)

// WebSocketProxy implements SignallingServerProxy using WebSocket
type WebSocketProxy struct {
	conn        *websocket.Conn
	nodeID      string
	nodeIDMu    sync.RWMutex
	receiveChan chan pkgframing.MessagePayload
	debug       bool
}

// NewWebSocketProxy creates a new WebSocket signalling proxy from an existing connection
func NewWebSocketProxy(conn *websocket.Conn, debug bool) *WebSocketProxy {
	proxy := &WebSocketProxy{
		conn:        conn,
		receiveChan: make(chan pkgframing.MessagePayload, 100),
		debug:       debug,
	}

	return proxy
}

// Send sends a message to the signalling server
func (p *WebSocketProxy) Send(ctx context.Context, msg pkgframing.MessagePayload) error {
	return p.conn.WriteJSON(msg)
}

// Receive returns a channel for receiving messages from the signalling server
func (p *WebSocketProxy) Receive() <-chan pkgframing.MessagePayload {
	return p.receiveChan
}

// readLoop reads messages from WebSocket and sends to receiveChan
func (p *WebSocketProxy) Run(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := p.conn.ReadMessage()
				if err != nil {
					log.Println("Read error:", err)
					return
				}

				if p.debug {
					log.Printf("Received: %s", message)
				}

				var payload pkgframing.MessagePayload
				if err := json.Unmarshal(message, &payload); err != nil {
					log.Printf("Failed to parse message: %v", err)
					continue
				}
				p.receiveChan <- payload
			}
		}
	}(ctx)
}

// FilteredSignallingProxy wraps a SignallingServerProxy and filters out pong messages
type FilteredSignallingProxy struct {
	proxy       pkgproxy.SignallingServerProxy
	receiveChan chan pkgframing.MessagePayload
	debug       bool
	wsConn      *websocket.Conn
}

// NewFilteredSignallingProxy creates a new filtered signalling proxy
func NewFilteredSignallingProxy(upstream pkgproxy.SignallingServerProxy, debug bool) *FilteredSignallingProxy {
	f := &FilteredSignallingProxy{
		proxy:       upstream,
		receiveChan: make(chan pkgframing.MessagePayload, 100),
		debug:       debug,
	}

	// Start goroutine to filter messages
	go f.filterLoop()

	return f
}

// filterLoop reads from the underlying proxy and filters out pong messages
func (f *FilteredSignallingProxy) filterLoop() {
	for msg := range f.proxy.Receive() {
		// Handle pong messages
		if msg.Echo != nil && msg.Echo.Direction == pkgconnreg.EchoDirectionS2C {
			if f.debug {
				rtt := time.Since(time.UnixMilli(int64(msg.Echo.Timestamp)))
				log.Printf("Pong received - RTT: %v, CorrelationID: %s, SeqID: %d",
					rtt, msg.Echo.CorrelationID, msg.Echo.SeqID)
			}
			continue
		}

		// Pass through non-pong messages
		f.receiveChan <- msg
	}
	close(f.receiveChan)
}

// Send sends a message through the underlying proxy
func (f *FilteredSignallingProxy) Send(ctx context.Context, msg pkgframing.MessagePayload) error {
	return f.proxy.Send(ctx, msg)
}

// Receive returns a channel that filters out pong messages
func (f *FilteredSignallingProxy) Receive() <-chan pkgframing.MessagePayload {
	return f.receiveChan
}
