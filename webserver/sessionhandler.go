package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

type SessionHandlerTimeProvider interface {
	GetCurrentTime() time.Time
	HasExpired(time.Time) bool
	StartTimer(f func())
}

type SessionHandlerTimeProviderImpl struct {
}

func (i SessionHandlerTimeProviderImpl) GetCurrentTime() time.Time {
	return time.Now()
}

func (i SessionHandlerTimeProviderImpl) HasExpired(t time.Time) bool {
	return t.AddDate(1, 0, 0).Before(time.Now())
}

func (i SessionHandlerTimeProviderImpl) StartTimer(f func()) {
	time.AfterFunc(time.Duration(1)*time.Hour, f)
}

type Session struct {
	UserName string
	UserType datacontract.UserType

	CreationTime time.Time
}

var sessionHandlerInstance *SessionHandler
var sessionHandlerOnce sync.Once

// SessionHandler manages the client sessions
type SessionHandler struct {
	syncObject sync.RWMutex

	sessions map[string]*Session

	timeProvider SessionHandlerTimeProvider
}

func GetSessionHandlerInstance() *SessionHandler {
	sessionHandlerOnce.Do(func() {
		sessionHandlerInstance = &SessionHandler{
			sessions:     make(map[string]*Session),
			timeProvider: &SessionHandlerTimeProviderImpl{},
		}

		go sessionHandlerInstance.CollectGarbage()
	})

	return sessionHandlerInstance
}

func createSessionID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (h *SessionHandler) CreateSession(userName string, userType datacontract.UserType) string {
	h.syncObject.Lock()
	defer h.syncObject.Unlock()

	sessionID := ""
	for id, session := range h.sessions {
		if userName == session.UserName {
			sessionID = id
			break
		}
	}

	if sessionID == "" {
		sessionID = createSessionID()
		h.sessions[sessionID] = &Session{
			UserName:     userName,
			UserType:     userType,
			CreationTime: h.timeProvider.GetCurrentTime(),
		}
	} else {
		h.sessions[sessionID].CreationTime = h.timeProvider.GetCurrentTime()
	}

	return sessionID
}

func (h *SessionHandler) RemoveSession(sessionID string) {
	h.syncObject.Lock()
	defer h.syncObject.Unlock()

	delete(h.sessions, sessionID)
}

func (h *SessionHandler) CollectGarbage() {
	h.syncObject.Lock()
	defer h.syncObject.Unlock()

	var sessionsToDelete []string
	for id, session := range h.sessions {
		if h.timeProvider.HasExpired(session.CreationTime) {
			sessionsToDelete = append(sessionsToDelete, id)
		}
	}

	for _, id := range sessionsToDelete {
		delete(h.sessions, id)
	}

	h.timeProvider.StartTimer(func() { h.CollectGarbage() })
}

func (h *SessionHandler) GetUserType(sessionID string) (hasSession bool, userType datacontract.UserType) {
	h.syncObject.RLock()
	defer h.syncObject.RUnlock()

	session, ok := h.sessions[sessionID]
	if !ok {
		return false, 0
	}

	return true, session.UserType
}
