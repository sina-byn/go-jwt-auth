package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sina-byn/go-jwt-auth/pkg/auth"
	"github.com/sina-byn/go-jwt-auth/pkg/db"
	"github.com/sina-byn/go-jwt-auth/pkg/user"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Failed to load environment variables from \".env\"")
	}

	db.Connect()
	defer db.DB.Close()

	r := gin.Default()

	user.RegisterRoutes(r)
	auth.RegisterRoutes(r)

	r.Run(":8080")
}
