package handler

import (
	"encoding/json"
	"net/http"
	"sync"
)

type CounterHandler struct {
	// key is session identifier, value is integer
	state sync.Map
}

func (cnt *CounterHandler) increaseCnt(key string) {
	// use a compare-and-set to update the per-session counter by one
}

func (cnt *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionId := ctx.Value(CtxSessionKeySessionId)
	if sessionId == nil {
		json.NewEncoder(w).Encode(&CounterHandlerResponse{Err: "No session identifier found"})
	}

	sessionIdStr, ok := sessionId.(string)
	if !ok {
		json.NewEncoder(w).Encode(&CounterHandlerResponse{Err: "No valid session identifier found"})
	}

	currentVal := 0
	if v, ok := cnt.state.Load(sessionIdStr); ok {
		currentVal = v.(int)
	}
	defer cnt.increaseCnt(sessionIdStr)

	json.NewEncoder(w).Encode(&CounterHandlerResponse{
		SessionId: sessionIdStr,
		Count:     currentVal,
	})
}

type CounterHandlerResponse struct {
	Err       string `json:"err,omitempty"`
	SessionId string `json:"session_id"`
	Count     int    `json:"count"`
}
