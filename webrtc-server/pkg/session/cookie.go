package session

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// This is a cookie-based session manager
// an implementation of handler.SessionManager

type CookieSessionManager struct {
	sessionStore   sync.Map
	CookieDomain   string
	CookieSameSite string
	CookieSecure   string
}

func (sessMngr *CookieSessionManager) getRandomSessionId() string {
	return uuid.NewString()
}

// If no associated is found with such request, returns an empty string
// otherwise returns a session identifier (which is opaque to the user)
func (sessMngr *CookieSessionManager) GetSessionId(ctx context.Context, r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}

	sessionId := cookie.Value
	// Check if the session exists in the store
	_, exists := sessMngr.sessionStore.Load(sessionId)
	if !exists {
		return ""
	}

	return sessionId
}

// Associate the request with a session, maybe alter the response when needed.
func (sessMngr *CookieSessionManager) CreateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	sessionId := sessMngr.getRandomSessionId()

	// Store the session
	sessMngr.sessionStore.Store(sessionId, struct{}{})

	// Set the cookie on the response
	cookieObj := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
	}
	if domain := sessMngr.CookieDomain; domain != "" {
		cookieObj.Domain = domain
	}
	if sameSite := sessMngr.CookieSameSite; sameSite != "" {
		switch sameSite {
		case "None":
			cookieObj.SameSite = http.SameSiteNoneMode
		case "Lax":
			cookieObj.SameSite = http.SameSiteLaxMode
		case "Strict":
			cookieObj.SameSite = http.SameSiteStrictMode
		default:
			log.Println("Unknown cookie sameSite setting", sameSite)
			cookieObj.SameSite = http.SameSiteDefaultMode
		}
	}
	if secure := sessMngr.CookieSecure; secure != "" {
		if secure == "true" {
			cookieObj.Secure = true
		} else if secure == "false" {
			cookieObj.Secure = false
		} else {
			log.Println("Unknown cookie secure setting", secure)
		}
	}

	// To understand why and when such options are usefull,
	// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Cookies

	http.SetCookie(w, cookieObj)

	return sessionId
}
