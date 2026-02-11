package main

import (
	"log"
	"net/http"
	"time"

	pkgconnreg "example.com/webrtcserver/pkg/connreg"
	pkghandler "example.com/webrtcserver/pkg/handler"
	pkgsafemap "example.com/webrtcserver/pkg/safemap"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	listenAddr := ":3001"
	wsTimeout := 10 * time.Second
	wsPath := "/ws"
	sm := pkgsafemap.NewSafeMap()
	cr := pkgconnreg.NewConnRegistry(sm)
	wsHandler := pkghandler.NewWebsocketHandler(&upgrader, cr, wsTimeout)

	var connsHandler http.Handler = pkghandler.NewConnsHandler(cr)
	connsHandler = pkghandler.WithCORSAllowAny(connsHandler)

	mux := http.NewServeMux()
	mux.Handle(wsPath, wsHandler)

	mux.Handle("/conns", connsHandler)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}
	log.Printf("Starting server on %s", listenAddr)
	server.ListenAndServe()
}
