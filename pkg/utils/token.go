package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func GenerateTokenPair(userId int64, email string) (*TokenPair, error) {
	jwtAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	if jwtAccessSecret == "" || jwtRefreshSecret == "" {
		return nil, errors.New("JWT secrets are not set in environment variables")
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	})

	signedAccessToken, err := accessToken.SignedString([]byte(jwtAccessSecret))

	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(jwtRefreshSecret))

	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
