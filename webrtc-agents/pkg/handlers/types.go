package handlers

import (
	"context"
	pkgwsrunner "webrtc-agents/pkg/ws_runner"
)

type GenericWebRTCHandler interface {
	Run(ctx context.Context, runner pkgwsrunner.WebSocketSignallingSessionRunner)
}
