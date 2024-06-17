package main

import (
	"log"

	"github.com/sina-byn/go-jwt-auth/pkg/auth"
	"github.com/sina-byn/go-jwt-auth/pkg/db"
	"github.com/sina-byn/go-jwt-auth/pkg/user"
	"github.com/sina-byn/go-jwt-auth/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load environment variables from \".env\"")
	}

	db.Connect()
	defer db.DB.Close()

	utils.InitTokenBlackList()

	r := gin.Default()

	user.RegisterRoutes(r)
	auth.RegisterRoutes(r)

	r.Run(":8080")
}
