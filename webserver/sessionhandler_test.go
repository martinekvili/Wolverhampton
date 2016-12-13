package main

import (
	"testing"
	"time"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

type SessionHandlerTimeProviderTestImpl struct {
	Expired bool
}

func (i SessionHandlerTimeProviderTestImpl) GetCurrentTime() time.Time {
	return time.Now()
}

func (i SessionHandlerTimeProviderTestImpl) HasExpired(t time.Time) bool {
	return i.Expired
}

func (i SessionHandlerTimeProviderTestImpl) StartTimer(f func()) {
	return
}

func setupSut() (*SessionHandlerTimeProviderTestImpl, SessionHandler) {
	mock := SessionHandlerTimeProviderTestImpl{
		Expired: false,
	}

	sut := SessionHandler{
		sessions:     make(map[string]*Session),
		timeProvider: &mock,
	}

	return &mock, sut
}

func TestSessionHandlerSimple(t *testing.T) {
	_, sut := setupSut()

	sessionID := sut.CreateSession("test", datacontract.Student)

	if sessionID == "" {
		t.Error("Did not get an actual sessionID.")
	}

	hasSession, userType := sut.GetUserType(sessionID)

	if !hasSession {
		t.Error("Session has not been created.")
	}

	if userType != datacontract.Student {
		t.Errorf("Expected user type was '%v', got '%v'.", datacontract.Student, userType)
	}

	sut.RemoveSession(sessionID)

	hasSession, _ = sut.GetUserType(sessionID)

	if hasSession {
		t.Error("Session should have been deleted.")
	}
}

func TestSessionHandlerSameUser(t *testing.T) {
	_, sut := setupSut()

	sessionID := sut.CreateSession("test", datacontract.Teacher)

	sessionID2 := sut.CreateSession("test", datacontract.Teacher)

	if sessionID != sessionID2 {
		t.Error("The two session IDs should be the same, but they are not.")
	}
}

func TestSessionHandlerExpiration(t *testing.T) {
	mock, sut := setupSut()

	sessionID := sut.CreateSession("test", datacontract.Admin)
	sut.CollectGarbage()

	hasSession, _ := sut.GetUserType(sessionID)
	if !hasSession {
		t.Error("Session has not been created.")
	}

	mock.Expired = true
	sut.CollectGarbage()

	hasSession, _ = sut.GetUserType(sessionID)
	if hasSession {
		t.Error("Session has not been removed.")
	}
}
