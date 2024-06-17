package middlewares

import (
	"net/http"
	"strings"

	"github.com/sina-byn/go-jwt-auth/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	token := strings.TrimSpace(strings.Split(authHeader, " ")[1])

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	tokenClaims, err := utils.VerifyToken(token, "access")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	c.Set("userId", tokenClaims.UserId)
	c.Set("email", tokenClaims.Email)
	c.Next()
}
