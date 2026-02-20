package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	pkghandlers "webrtc-agents/pkg/handlers"

	pkgframing "example.com/webrtcserver/pkg/framing"

	"github.com/alecthomas/kong"
	"github.com/gorilla/websocket"
)

var cli struct {
	WsServer              string        `name:"ws-server" help:"WebSocket server URL" default:"ws://localhost:3001/ws"`
	NodeName              string        `name:"node-name" help:"Node name for registration" default:"webrtc-agent-1"`
	PingPeriodSeconds     int           `name:"ping-period-seconds" help:"Ping period in seconds" default:"5"`
	Debug                 bool          `name:"debug" help:"Show ping/pong messages in logs for debugging purposes"`
	ICEServer             []string      `name:"ice-server" help:"To specify the ICE servers, might be specify multiple times" default:"stun:stun.l.google.com:19302"`
	ReconnectOnDisconnect bool          `name:"reconnect-on-disconnect" help:"Reconnect on WebSocket disconnect"`
	ReconnectDelay        time.Duration `name:"reconnect-delay" help:"Delay between reconnect attempts" default:"3s"`
}

type WebSocketRunner struct {
	ReconnectDelay        time.Duration
	ReconnectOnDisconnect bool
	PingIntv              time.Duration
	Debug                 bool
}

func (runner *WebSocketRunner) Run(ctx context.Context, u url.URL, txChannel <-chan pkgframing.MessagePayload) chan pkgframing.MessagePayload {
	outputDataCh := make(chan pkgframing.MessagePayload)
	go func(ctx context.Context) {
		defer close(outputDataCh)
		for {
			log.Printf("Connecting to %s", u.String())

			// Establish WebSocket connection
			wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("Failed to dial:", err)
			}
			defer wsConn.Close()

			registerer := &pkghandlers.WebSocketRegisterer{}
			if err := registerer.Register(wsConn, cli.NodeName); err != nil {
				log.Fatal("Failed to send registration message:", err)
			}

			log.Printf("Sent registration message for node: %s", cli.NodeName)

			wsPinger := &pkghandlers.WebSocketPinger{
				Intv:  runner.PingIntv,
				Debug: runner.Debug,
			}
			dataCh, errCh := wsPinger.StartPingLoop(ctx, wsConn, txChannel)
			log.Println("Ping/pong loop started")

			go func() {
				for item := range dataCh {
					outputDataCh <- item
				}
			}()

			err, ok := <-errCh
			if ok && err != nil {
				log.Printf("Error on ws connection: %+v", err)
				if !runner.ReconnectOnDisconnect {
					return
				}
				break
			}

			log.Printf("Reconnecting to %s in %s", u.String(), runner.ReconnectDelay.String())
			<-time.After(runner.ReconnectDelay)
		}
	}(ctx)

	return outputDataCh
}

func main() {
	kong.Parse(&cli)

	pingPeriod := time.Duration(cli.PingPeriodSeconds) * time.Second

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Parse WebSocket URL
	u, err := url.Parse(cli.WsServer)
	if err != nil {
		log.Fatal("Failed to parse WebSocket URL:", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := &WebSocketRunner{
		PingIntv:              pingPeriod,
		Debug:                 cli.Debug,
		ReconnectOnDisconnect: cli.ReconnectOnDisconnect,
		ReconnectDelay:        cli.ReconnectDelay,
	}
	signallingTxChannel := make(chan pkgframing.MessagePayload)
	signallingDataCh := runner.Run(ctx, *u, signallingTxChannel)

	// Create WebRTC handler
	webrtcHandler := pkghandlers.NewWebRTCHandler(cli.ICEServer, cli.Debug, signallingTxChannel, signallingDataCh)

	webrtcHandler.Run(ctx)

	sigsCh := make(chan os.Signal, 1)
	signal.Notify(sigsCh, syscall.SIGINT)
	signal := <-sigsCh
	log.Printf("Received signal %+v, exitting", signal.String())
}
