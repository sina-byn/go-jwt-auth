package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	userGroup := r.Group("/auth")

	userGroup.POST("/login", getUsers)
	userGroup.POST("/logout")

	return userGroup
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, "done")
	fmt.Println("this is executed")
}
