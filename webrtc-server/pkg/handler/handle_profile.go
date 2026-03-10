package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	pkglogin "example.com/webrtcserver/pkg/models/login"
	pkguser "example.com/webrtcserver/pkg/models/user"
)

type ProfileHandler struct {
	// Get the user object from userId
	UserManager pkguser.UserManager

	// Check if current session has logged in
	UserSessionManager pkglogin.UserSessionManager
}

type ProfileResponse struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarURL,omitempty"`
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	if specifiedUsername := r.URL.Query().Get(QueryParamUsername); specifiedUsername != "" {
		userObj, err := h.UserManager.GetUserByUsername(ctx, specifiedUsername)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't access internal user store"})
			return
		}
		if userObj == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrResponse{Err: "User is not found"})
			return
		}
		err = json.NewEncoder(w).Encode(&ProfileResponse{
			Username:    userObj.Username,
			DisplayName: userObj.DisplayName,
			AvatarURL:   userObj.AvatarURL,
		})
		if err != nil {
			log.New(os.Stderr, "", 0).Printf("Cant format response: %v", err)
		}
		return
	}

	sessId := ctx.Value(CtxSessionKeySessionId)
	if sessId == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "No session Id is found"})
		return
	}

	userId, err := h.UserSessionManager.GetUserIdBySessionId(ctx, sessId.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't determine if you has logged in"})
		return
	}

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Unauthorized"})
		return
	}

	userObj, err := h.UserManager.GetUserById(ctx, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't access internal user store"})
		return
	}

	if userObj == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "User didn't found"})
		return
	}

	err = json.NewEncoder(w).Encode(&ProfileResponse{
		Username:    userObj.Username,
		DisplayName: userObj.DisplayName,
		AvatarURL:   userObj.AvatarURL,
	})
	if err != nil {
		log.New(os.Stderr, "", 0).Printf("Cant format response: %v", err)
	}
}

type ProfileStatusHandler struct {
	// Get the user object from userId
	UserManager pkguser.UserManager

	// Check if current session has logged in
	UserSessionManager pkglogin.UserSessionManager
}

type ProfileStatusResponse struct {
	LoggedIn bool `json:"logged_in"`
}

type ProfileAvatarHandler struct {
	// Get the user object from userId
	UserManager pkguser.UserManager

	// Check if current session has logged in
	UserSessionManager pkglogin.UserSessionManager
}

func (h *ProfileAvatarHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	if specifiedUsername := r.URL.Query().Get(QueryParamUsername); specifiedUsername != "" {
		userObj, err := h.UserManager.GetUserByUsername(ctx, specifiedUsername)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't access internal user store"})
			return
		}
		if userObj == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrResponse{Err: "User is not found"})
			return
		}
		h.serveUserAvatar(w, userObj)
		return
	}

	sessId := ctx.Value(CtxSessionKeySessionId)
	if sessId == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "No session Id is found"})
		return
	}

	userId, err := h.UserSessionManager.GetUserIdBySessionId(ctx, sessId.(string))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't determine if you has logged in"})
		return
	}

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Unauthorized"})
		return
	}

	userObj, err := h.UserManager.GetUserById(ctx, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Can't access internal user store"})
		return
	}

	h.serveUserAvatar(w, userObj)
}

func (h *ProfileAvatarHandler) serveUserAvatar(w http.ResponseWriter, userObj *pkguser.User) {
	if userObj == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "User didn't found"})
		return
	}

	if userObj.AvatarURL == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "No avatar URL found"})
		return
	}

	resp, err := http.Get(userObj.AvatarURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Failed to fetch avatar"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(&ErrResponse{Err: "Upstream returned error"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.New(os.Stderr, "", 0).Printf("Cant write avatar response: %v", err)
	}
}

func (h *ProfileStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if sessId := ctx.Value(CtxSessionKeySessionId); sessId != nil {
		if userId, err := h.UserSessionManager.GetUserIdBySessionId(ctx, sessId.(string)); err == nil && userId != "" {
			if u, err := h.UserManager.GetUserById(ctx, userId); err == nil && u != nil {
				json.NewEncoder(w).Encode(&ProfileStatusResponse{
					LoggedIn: true,
				})
				return
			}
		}
	}
	json.NewEncoder(w).Encode(&ProfileStatusResponse{LoggedIn: false})
}
