package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	pkguser "example.com/webrtcserver/pkg/models/user"
	pkgmyjwt "example.com/webrtcserver/pkg/my_jwt"
)

type BotsManagementHandler struct {
	UserManager pkguser.UserManager
	JWTManager  pkgmyjwt.JWTManager
}

type BotCreationPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

func (h *BotsManagementHandler) handleAddBot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&ErrResponse{Err: pkguser.ErrUsernameDuplicated.Error()})
}

func (h *BotsManagementHandler) handleDeleteBot(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&ErrResponse{Err: "not implemented yet"})
}

func (h *BotsManagementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/add") && r.Method == http.MethodPost {
		h.handleAddBot(w, r)
		return
	} else if strings.HasSuffix(r.URL.Path, "/delete") && r.Method == http.MethodDelete {
		h.handleDeleteBot(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&ErrResponse{Err: "no handlers matched the request."})
}
