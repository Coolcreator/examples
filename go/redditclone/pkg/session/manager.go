package session

import (
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionsManager struct {
	data  map[string]*Session
	mu    *sync.RWMutex
	token *jwt.Token
}

func NewSessionsMem() *SessionsManager {
	return &SessionsManager{
		data:  make(map[string]*Session, 10),
		mu:    &sync.RWMutex{},
		token: nil,
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, ErrNoAuth
	}

	sm.mu.RLock()
	sess, ok := sm.data[sessionCookie.Value]
	sm.mu.RUnlock()

	if !ok {
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32, userLogin string) (*Session, *jwt.Token, error) {
	sess := NewSession(userID, userLogin)

	sm.mu.Lock()
	sm.data[sess.ID] = sess
	sm.mu.Unlock()

	sm.token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": jwt.MapClaims{"username": sess.Login, "id": sess.ID},
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return sess, sm.token, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	sm.mu.Lock()
	delete(sm.data, sess.ID)
	sm.mu.Unlock()

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
