package api

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

const headwayTriggerKey = "headway_trigger"

// JobTokenAuth validates Authorization: Bearer tokens for automated jobs (e.g. GitHub Actions).
type JobTokenAuth struct {
	token string
}

func JobTokenAuthFromEnv() (*JobTokenAuth, bool) {
	token := strings.TrimSpace(os.Getenv("HEADWAY_JOB_TOKEN"))
	if token == "" {
		return nil, false
	}
	return &JobTokenAuth{token: token}, true
}

func (j *JobTokenAuth) ValidBearer(header string) bool {
	if j == nil || j.token == "" {
		return false
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return false
	}
	got := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if got == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(j.token)) == 1
}

// headwayAuthMiddleware accepts either an admin session cookie or the job bearer token.
func (a *API) headwayAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.adminAuth != nil && a.adminAuth.ValidateToken(a.adminAuth.sessionToken(c)) {
			c.Set(headwayTriggerKey, business.HeadwayTriggerAdmin)
			c.Next()
			return
		}
		if a.jobTokenAuth != nil && a.jobTokenAuth.ValidBearer(c.GetHeader("Authorization")) {
			c.Set(headwayTriggerKey, business.HeadwayTriggerCron)
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func headwayTriggerFromContext(c *gin.Context) business.HeadwayJobTrigger {
	if v, ok := c.Get(headwayTriggerKey); ok {
		if t, ok := v.(business.HeadwayJobTrigger); ok {
			return t
		}
	}
	return business.HeadwayTriggerAPI
}
