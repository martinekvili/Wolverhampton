package main

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

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
}

func GetSessionHandlerInstance() *SessionHandler {
	sessionHandlerOnce.Do(func() {
		sessionHandlerInstance = &SessionHandler{
			sessions: make(map[string]*Session),
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
			CreationTime: time.Now(),
		}
	} else {
		h.sessions[sessionID].CreationTime = time.Now()
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
		if session.CreationTime.AddDate(1, 0, 0).Before(time.Now()) {
			sessionsToDelete = append(sessionsToDelete, id)
		}
	}

	for _, id := range sessionsToDelete {
		delete(h.sessions, id)
	}

	time.AfterFunc(time.Duration(1)*time.Hour, func() { h.CollectGarbage() })
}
