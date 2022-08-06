package session

import (
	"database/sql"
	"net/http"
	"time"
)

type SessionsManager struct {
	data *sql.DB
}

func NewSessionsRepo(db *sql.DB) *SessionsManager {
	return &SessionsManager{data: db}
}
func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	sess := &Session{}
	errD := sm.data.
		QueryRow("SELECT data, userID FROM sessions WHERE data= ?", sessionCookie.Value).
		Scan(&sess.ID, &sess.UserID)
	if errD != nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID int64) (*Session, error) {
	sess := NewSession(userID)

	_, err := sm.data.Exec(
		"INSERT INTO  sessions (`data`, `userID`) VALUES (?, ?)",
		sess.ID,
		sess.UserID,
	)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	_, errD := sm.data.Exec(
		"DELETE FROM sessions WHERE data = ?",
		sess.ID,
	)
	if errD != nil {
		return errD
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
