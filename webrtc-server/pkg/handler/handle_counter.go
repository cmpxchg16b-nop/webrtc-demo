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

// returns the pre-increment value
func (cnt *CounterHandler) increaseCnt(key string) int {
	// use a compare-and-swap loop to handle concurrent updates
	currentVal := 0
	for {
		// CompareAndSwap failed, re-read the current value and retry
		if v, loaded := cnt.state.LoadOrStore(key, currentVal); loaded {
			currentVal = v.(int)
		}

		if cnt.state.CompareAndSwap(key, currentVal, currentVal+1) {
			break
		}
	}

	return currentVal
}

func (cnt *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	sessionId := ctx.Value(CtxSessionKeySessionId)
	if sessionId == nil {
		json.NewEncoder(w).Encode(&CounterHandlerResponse{Err: "No session identifier found"})
		return
	}

	sessionIdStr, ok := sessionId.(string)
	if !ok {
		json.NewEncoder(w).Encode(&CounterHandlerResponse{Err: "No valid session identifier found"})
		return
	}

	currentVal := cnt.increaseCnt(sessionIdStr)

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
