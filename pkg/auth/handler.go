package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sina-byn/go-jwt-auth/pkg/user"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	authGroup := r.Group("/auth")

	authGroup.POST("/login", loginHandler)
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

func refreshHandler(c *gin.Context) {
	var bodyStruct struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	err := c.ShouldBindJSON(&bodyStruct)

	if err != nil {
		c.JSON(http.StatusBadRequest, "Could not parse request data")
		return
	}

	refreshedTokenPair, err := Refresh(bodyStruct.RefreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh tokens"})
		return
	}

	c.JSON(http.StatusOK, refreshedTokenPair)
}
