package main

import (
	"github.com/sina-byn/go-jwt-auth/pkg/db"
	"github.com/sina-byn/go-jwt-auth/pkg/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.DB.Close()

	r := gin.Default()

	user.RegisterRoutes(r)

	r.Run(":8080")
}
