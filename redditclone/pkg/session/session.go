package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
)

type Session struct {
	ID     string
	UserID int64
}

func NewSession(userID int64) *Session {
	randID := make([]byte, 16)
	rand.Read(randID)

	return &Session{
		ID:     fmt.Sprintf("%x", randID),
		UserID: userID,
	}
}

var (
	ErrNoAuth = errors.New("No session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

//go:generate mockgen -source=session.go -destination=repo_mock.go -package=post SessRepo
type SessRepo interface {
	Create(w http.ResponseWriter, userID int64) (*Session, error)
	DestroyCurrent(w http.ResponseWriter, r *http.Request) error
	Check(r *http.Request) (*Session, error)
}
