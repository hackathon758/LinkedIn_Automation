package auth

import (
	"time"

	"linkedin-automation/database"
)

// SessionManager manages LinkedIn sessions
type SessionManager struct {
	db         *database.DB
	isLoggedIn bool
	lastCheck  time.Time
}

// NewSessionManager creates a new session manager
func NewSessionManager(db *database.DB) *SessionManager {
	return &SessionManager{
		db: db,
	}
}

// HasValidSession checks if there's a valid stored session
func (sm *SessionManager) HasValidSession() bool {
	cookies, err := sm.db.GetCookies()
	if err != nil || len(cookies) == 0 {
		return false
	}

	// Check if any essential cookies are present and not expired
	essentialCookies := []string{"li_at", "JSESSIONID"}
	for _, essential := range essentialCookies {
		found := false
		for _, cookie := range cookies {
			if cookie.Name == essential && cookie.ExpiresAt.After(time.Now()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// SetLoggedIn updates the login state
func (sm *SessionManager) SetLoggedIn(loggedIn bool) {
	sm.isLoggedIn = loggedIn
	sm.lastCheck = time.Now()
}

// IsLoggedIn returns the current login state
func (sm *SessionManager) IsLoggedIn() bool {
	return sm.isLoggedIn
}

// ClearSession clears the stored session
func (sm *SessionManager) ClearSession() error {
	sm.isLoggedIn = false
	return sm.db.ClearCookies()
}

// GetSessionAge returns how old the current session is
func (sm *SessionManager) GetSessionAge() time.Duration {
	cookies, err := sm.db.GetCookies()
	if err != nil || len(cookies) == 0 {
		return 0
	}

	// Find the oldest cookie
	var oldest time.Time
	for _, cookie := range cookies {
		if oldest.IsZero() || cookie.CreatedAt.Before(oldest) {
			oldest = cookie.CreatedAt
		}
	}

	return time.Since(oldest)
}

// NeedsRefresh checks if the session should be refreshed
func (sm *SessionManager) NeedsRefresh() bool {
	// Refresh if session is older than 12 hours
	return sm.GetSessionAge() > 12*time.Hour
}
