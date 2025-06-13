package sessions

import (
	"errors"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	UserId       int
	SessionToken string
	ExpireTime   time.Time
}

const sessionExpireTime = 100 * time.Minute

var Sessions map[string]Session = map[string]Session{}

// Create new session and login/register
func NewSession(userId int, w http.ResponseWriter) error {
	token, err := uuid.NewV4()
	if err != nil {
		return errors.New("error creating identification number")
	}
	session := Session{
		UserId:       userId,
		SessionToken: token.String(),
		ExpireTime:   time.Now().Add(sessionExpireTime),
	}
	Sessions[token.String()] = session
	session.setCookie(w)

	return nil
}

// Validate session existence and get userId
// Used in every handler
func Validate(r *http.Request) (int, error) {
	cookie, err := getCookie(r)
	if err != nil {
		return 0, err
	}
	if IsExpired(cookie) {
		return -1, errors.New("session is expired")
	}
	key, ok := Sessions[cookie.Value]
	if !ok {
		return -1, errors.New("there is no such session")
	}
	return key.UserId, nil
}

// Close session when logout or session expired
func CloseSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := getCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	key, ok := Sessions[cookie.Value]
	if !ok {
		log.Printf("Error closing session: session is already closed\n")
		return
	}
	delete(Sessions, key.SessionToken)
	newCookie := http.Cookie{
		Name:   "TOKEN",
		MaxAge: -1,
	}
	http.SetCookie(w, &newCookie)
}

// Check if session is expired
func IsExpired(cookie *http.Cookie) bool {
	key, ok := Sessions[cookie.Value]
	if !ok { // Check if session existence
		return true
	}
	if key.ExpireTime.Before(time.Now()) {
		return true
	}
	// Update the session expire time
	key.ExpireTime = time.Now().Add(sessionExpireTime) //1 year
	return false
}

// Get cookie from HEADER of browser for validation
func getCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("TOKEN")
	if err != nil {
		log.Printf("Error getting cookie: %v\n", err)
		return nil, err
	}
	return cookie, nil
}

// Set new Cookie
func (s *Session) setCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    "TOKEN",
		Value:   s.SessionToken,
		Expires: time.Now().Add(math.MaxInt),
		MaxAge:  math.MaxInt,
		Path:    "/",
	}

	if err := cookie.Valid(); err != nil {
		log.Println(err)
	}
	http.SetCookie(w, &cookie)
}
