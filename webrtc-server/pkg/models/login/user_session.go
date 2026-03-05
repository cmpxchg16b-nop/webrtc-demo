package login

import (
	"context"
	"sync"
)

type UserSessionManager interface {
	// Login is to associate a session with a registered user
	LogIn(ctx context.Context, userId string, sessionId string) error

	// Log out is to de-ssociate a session with the user bound
	LogOut(ctx context.Context, sessionId string) error

	// Returns an empty string if the user hasn't log in, otherwise returns an non-empty string
	// Callers should check error first
	GetUserIdBySessionId(ctx context.Context, sessionId string) (string, error)
}

type MemoryUserSessionManager struct {
	store sync.Map
}

// LogIn associates a session with a registered user
func (m *MemoryUserSessionManager) LogIn(ctx context.Context, userId string, sessionId string) error {
	m.store.Store(sessionId, userId)
	return nil
}

// LogOut de-associates a session with the user bound
func (m *MemoryUserSessionManager) LogOut(ctx context.Context, sessionId string) error {
	m.store.Delete(sessionId)
	return nil
}

// GetUserIdBySessionId returns the user ID associated with the session
// Returns an empty string if the user hasn't logged in
func (m *MemoryUserSessionManager) GetUserIdBySessionId(ctx context.Context, sessionId string) (string, error) {
	if v, ok := m.store.Load(sessionId); ok {
		return v.(string), nil
	}
	return "", nil
}
