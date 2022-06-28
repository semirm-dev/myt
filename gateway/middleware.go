package gateway

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BasicAuth implementation ref: https://www.alexedwards.net/blog/basic-authentication-in-go
func BasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usr, pwd, ok := c.Request.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(usr))
			passwordHash := sha256.Sum256([]byte(pwd))
			expectedUsernameHash := sha256.Sum256([]byte(username))
			expectedPasswordHash := sha256.Sum256([]byte(password))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				c.Next()
				return
			}
		}

		c.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
