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
func (pinger *WebSocketPinger) StartPingLoop(ctx context.Context, wsConn *websocket.Conn, txChannel <-chan pkgframing.MessagePayload) (chan pkgframing.MessagePayload, chan error) {
	period := pinger.Intv
	debug := pinger.Debug

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	seqID := uint64(0)

	errCh := make(chan error, 1)
	dataCh := make(chan pkgframing.MessagePayload, 1)

	go func(ctx context.Context) {
		defer close(dataCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return
			case txItem := <-txChannel:
				if err := wsConn.WriteJSON(txItem); err != nil {
					log.Println("Failed to send message:", err)
					errCh <- err
					return
				}
				continue
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
					errCh <- err
					return
				}

				if debug {
					log.Printf("Sent ping - SeqID: %d, CorrelationID: ping-%d", seqID, seqID)
				}
			default:
				var msg pkgframing.MessagePayload
				if err := wsConn.ReadJSON(&msg); err != nil {
					log.Println("Failed to get message from ws connection:", err)
					errCh <- err
					return
				}
				if msg.Echo != nil {
					if msg.Echo.Direction == pkgconnreg.EchoDirectionS2C && pinger.Debug {
						rtt := time.Since(time.UnixMilli(int64(msg.Echo.Timestamp)))
						log.Printf("Pong received - RTT: %v, CorrelationID: %s, SeqID: %d",
							rtt, msg.Echo.CorrelationID, msg.Echo.SeqID)
					}
					continue
				}
				dataCh <- msg
			}
		}
	}(ctx)

	return dataCh, errCh
}
