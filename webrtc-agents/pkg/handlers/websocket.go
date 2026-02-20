package handlers

import (
	"context"
	"fmt"
	"log"

	"time"

	pkgconnreg "example.com/webrtcserver/pkg/connreg"
	pkgframing "example.com/webrtcserver/pkg/framing"
	"github.com/gorilla/websocket"
)

type WebSocketRegisterer struct{}

// Send registration message
func (reg *WebSocketRegisterer) Register(wsConn *websocket.Conn, nodeName string) error {
	registerMsg := pkgframing.MessagePayload{
		Register: &pkgconnreg.RegisterPayload{
			NodeName: nodeName,
		},
	}

	return wsConn.WriteJSON(registerMsg)
}

type WebSocketPinger struct {
	Intv  time.Duration
	Debug bool
}

// Start ping goroutine to maintain WebSocket connection
func (pinger *WebSocketPinger) StartPingLoop(ctx context.Context, wsConn *websocket.Conn) {
	period := pinger.Intv
	debug := pinger.Debug

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	seqID := uint64(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			seqID++
			now := uint64(time.Now().UnixMilli())
			pingMsg := pkgframing.MessagePayload{
				Echo: &pkgconnreg.EchoPayload{
					Direction:     pkgconnreg.EchoDirectionC2S,
					CorrelationID: fmt.Sprintf("ping-%d", seqID),
					Timestamp:     now,
					SeqID:         seqID,
				},
			}

			if err := wsConn.WriteJSON(pingMsg); err != nil {
				log.Println("Failed to send ping:", err)
				return
			}

			if debug {
				log.Printf("Sent ping - SeqID: %d, CorrelationID: ping-%d", seqID, seqID)
			}
		}
	}
}
