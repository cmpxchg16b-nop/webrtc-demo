package ws_proxy

import (
	"context"
	"errors"

	"github.com/gorilla/websocket"
)

type WebsocketWriteProxy struct {
	Conn         *websocket.Conn
	requestsChan chan chan WriteRequest
}

func NewWebsocketWriteProxy(conn *websocket.Conn) *WebsocketWriteProxy {
	return &WebsocketWriteProxy{
		Conn: conn,
	}
}

type WriteRequest struct {
	Data interface{}
	Err  chan error
}

func (ws_proxy *WebsocketWriteProxy) Run(ctx context.Context) {
	ws_proxy.requestsChan = make(chan chan WriteRequest)
	go func(ctx context.Context) {
		defer close(ws_proxy.requestsChan)

		for {

			requestC := make(chan WriteRequest)

			select {
			case <-ctx.Done():
				return
			case ws_proxy.requestsChan <- requestC:
				request, ok := <-requestC
				if !ok {
					continue
				}
				request.Err <- ws_proxy.Conn.WriteJSON(request.Data)
			}
		}
	}(ctx)
}

func (ws_proxy *WebsocketWriteProxy) WriteJSON(data interface{}) error {
	requestC, ok := <-ws_proxy.requestsChan
	if !ok {
		return errors.New("websocket proxy is already closed.")
	}

	request := WriteRequest{
		Data: data,
		Err:  make(chan error),
	}
	requestC <- request
	return <-request.Err
}
