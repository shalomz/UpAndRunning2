package lib

import (
	"crypto/rand"

	"github.com/op/go-logging"
	"net/http"

	"strings"
	"time"

	"encoding/base64"
	"io"
)

var cookieStorage map[string]http.Cookie

func InitSessionManagement() {
	logging.MustGetLogger("logger").Debug("Initializing Session-Management...")
	cookieStorage = make(map[string]http.Cookie)
}

func LoginAndGetCookie(username string) http.Cookie {
	// sessionValue: username + randomString
	randomId := getRandomId()
	sessionValue := strings.TrimSpace(username + ":" + randomId)

	// Build Cookie
	cookie := http.Cookie{Name: "session", Value: sessionValue, Path: "/", Expires: time.Now().AddDate(0, 0, 14), HttpOnly: true}

	// Save and return Cookie
	cookieStorage[username] = cookie
	return cookie
}

func IsLoggedIn(r *http.Request) bool {
	// Get Cookie from Request
	rCookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// Get data from received Cookie
	rCookieData := strings.Split(rCookie.Value, ":")

	// Get data from saved Cookie
	if _, ok := cookieStorage[rCookieData[0]]; !ok {
		return false
	}
	sCookie := cookieStorage[rCookieData[0]]
	sCookieData := strings.Split(sCookie.Value, ":")

	// Do not allow expired Cookies in Storage
	if sCookie.Expires.Before(time.Now()) {
		delete(cookieStorage, rCookieData[0])
		return false
	}

	// Check if the saved Cookie's randomId equals the received Cookie's randomId
	return rCookieData[1] == sCookieData[1]
}

func LogoutAndDestroyCookie(r *http.Request) http.Cookie {
	cookie, _ := r.Cookie("session")

	// Remove the saved Cookie
	delete(cookieStorage, strings.Split(cookie.Value, ":")[0])
	logging.MustGetLogger("logger").Info("Logout successful.")

	// Return useless Cookie
	return http.Cookie{Name: "session", Value: "", Path: "/", Expires: time.Now().AddDate(0, 0, -1), HttpOnly: true}
}

func getRandomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
