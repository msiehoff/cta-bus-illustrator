package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const adminSessionCookie = "admin_session"

type AdminAuth struct {
	username string
	password string
	secret   []byte
}

type adminSession struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func AdminAuthFromEnv() (*AdminAuth, bool) {
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	password := os.Getenv("ADMIN_PASSWORD")
	if username == "" || password == "" {
		return nil, false
	}

	secret := []byte(strings.TrimSpace(os.Getenv("ADMIN_SESSION_SECRET")))
	if len(secret) == 0 {
		sum := sha256.Sum256([]byte(username + ":" + password))
		secret = sum[:]
	}

	return &AdminAuth{
		username: username,
		password: password,
		secret:   secret,
	}, true
}

func (a *AdminAuth) Login(username, password string) (string, error) {
	if username != a.username || password != a.password {
		return "", errors.New("invalid credentials")
	}

	session := adminSession{
		Sub: username,
		Exp: time.Now().Add(24 * time.Hour).Unix(),
	}
	return a.signSession(session)
}

func (a *AdminAuth) signSession(session adminSession) (string, error) {
	payload, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, a.secret)
	mac.Write(payload)
	sig := mac.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(payload) + "." +
		base64.RawURLEncoding.EncodeToString(sig), nil
}

func (a *AdminAuth) ValidateToken(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, a.secret)
	mac.Write(payload)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return false
	}

	var session adminSession
	if err := json.Unmarshal(payload, &session); err != nil {
		return false
	}
	if session.Sub != a.username {
		return false
	}
	return time.Now().Unix() <= session.Exp
}

func (a *AdminAuth) setSessionCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(adminSessionCookie, token, 86400, "/", "", false, true)
}

func (a *AdminAuth) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(adminSessionCookie, "", -1, "/", "", false, true)
}

func (a *AdminAuth) sessionToken(c *gin.Context) string {
	token, err := c.Cookie(adminSessionCookie)
	if err != nil {
		return ""
	}
	return token
}

func (a *AdminAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "admin auth not configured"})
			return
		}
		if !a.ValidateToken(a.sessionToken(c)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func parseIntQuery(c *gin.Context, key string, defaultVal, max int) int {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return defaultVal
	}
	if n > max {
		return max
	}
	return n
}
