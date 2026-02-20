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
	NodeName              string        `name:"node-name" help:"Node name for registration" default:"EchoBot"`
	PingPeriod            time.Duration `name:"ping-period-seconds" help:"Ping period in seconds" default:"3s"`
	Debug                 bool          `name:"debug" help:"Show ping/pong messages in logs for debugging purposes"`
	ICEServer             []string      `name:"ice-server" help:"To specify the ICE servers, might be specify multiple times" default:"stun:stun.l.google.com:19302"`
	ReconnectOnDisconnect bool          `name:"reconnect-on-disconnect" help:"Reconnect on WebSocket disconnect"`
	ReconnectDelay        time.Duration `name:"reconnect-delay" help:"Delay between reconnect attempts" default:"3s"`
}

type WebSocketRunner struct {
	URL                   url.URL
	ReconnectDelay        time.Duration
	ReconnectOnDisconnect bool
	PingIntv              time.Duration
	Debug                 bool
}

func (runner *WebSocketRunner) Run(ctx context.Context, txChannel <-chan pkgframing.MessagePayload) chan pkgframing.MessagePayload {
	u := runner.URL
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
			log.Printf("Dialed to ws server %+v", wsConn.RemoteAddr().String())

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

	// Parse WebSocket URL
	u, err := url.Parse(cli.WsServer)
	if err != nil {
		log.Fatal("Failed to parse WebSocket URL:", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := &WebSocketRunner{
		URL:                   *u,
		PingIntv:              cli.PingPeriod,
		Debug:                 cli.Debug,
		ReconnectOnDisconnect: cli.ReconnectOnDisconnect,
		ReconnectDelay:        cli.ReconnectDelay,
	}

	// wsconn might be re-dial anytime, so the webrtc handler can't just send message directly to wsconn,
	// it send messages to this tx channel, which in turn forward messages to wsconn
	signallingTxChannel := make(chan pkgframing.MessagePayload)

	// webrtc handler get signalling channel messages from this channel, instead of reading wsconn directly,
	// the reason is the same as above.
	signallingDataCh := runner.Run(ctx, signallingTxChannel)

	webrtcHandler := pkghandlers.NewWebRTCHandler(cli.ICEServer, cli.Debug, signallingTxChannel, signallingDataCh)

	webrtcHandler.Run(ctx)

	sigsCh := make(chan os.Signal, 1)
	signal.Notify(sigsCh, syscall.SIGINT)
	signal := <-sigsCh
	log.Printf("Received signal %+v, exitting", signal.String())
}
