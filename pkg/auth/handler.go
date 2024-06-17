package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sina-byn/go-jwt-auth/pkg/blacklist"
	"github.com/sina-byn/go-jwt-auth/pkg/user"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	authGroup := r.Group("/auth")

	authGroup.POST("/login", loginHandler)
	authGroup.POST("/logout", logoutHandler)
	authGroup.POST("/refresh", refreshHandler)

	return authGroup
}

func loginHandler(c *gin.Context) {
	var user user.User

	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request data"})
		return
	}

	token, err := Login(user.Email, user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not authenticate user"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func logoutHandler(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	var RequestBody struct {
		RefreshToken string `json:"refreshToken"`
	}

	err := c.ShouldBindJSON(&RequestBody)

	if err != nil {
		RequestBody.RefreshToken = ""
	}

	if authHeader == "" && RequestBody.RefreshToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No tokens provided"})
		return
	}

	accessToken := strings.TrimSpace(strings.Split(authHeader, " ")[1])

	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token provided"})
		return
	}

	Logout(accessToken, RequestBody.RefreshToken)

	c.Status(http.StatusNoContent)
}

func refreshHandler(c *gin.Context) {
	var RequestBody struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	err := c.ShouldBindJSON(&RequestBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Could not parse request data")
		return
	}

	if blacklist.BlockedTokens.Blocked(RequestBody.RefreshToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Revoked token"})
		return
	}

	refreshedTokenPair, err := Refresh(RequestBody.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh tokens"})
		return
	}

	c.JSON(http.StatusOK, refreshedTokenPair)
}
