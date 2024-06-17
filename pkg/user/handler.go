package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sina-byn/go-jwt-auth/pkg/middlewares"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	userGroup := r.Group("/user")

	userGroup.GET("/:email", middlewares.Authenticate, getUserByEmailHandler)
	userGroup.POST("/", createUserHandler)
	userGroup.PUT("/:id", middlewares.Authenticate, updateUserHandler)
	userGroup.DELETE("/:id", middlewares.Authenticate, deleteUserHandler)

	return userGroup
}

func getUserByEmailHandler(c *gin.Context) {
	email := c.Param("email")

	user, err := GetUserByEmail(email)

	if user == nil && err == nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := UserResponse{
		Id:       user.Id,
		Email:    user.Email,
		Fullname: user.Fullname,
	}

	c.JSON(http.StatusOK, userResponse)
}

func createUserHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := CreateUser(user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Id = *userId

	c.JSON(http.StatusCreated, user)
}

func updateUserHandler(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = UpdateUser(intId, &user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func deleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	intId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	err = DeleteUser(intId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
