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

type TokenClaims struct {
	UserId int64
	Email  string
}

func GenerateTokenPair(userId int64, email string) (*TokenPair, error) {
	var jwtAccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	var jwtRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")

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
		"email":  email,
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

func VerifyToken(token, tokenType string) (*TokenClaims, error) {
	var jwtAccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	var jwtRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		secret := jwtAccessSecret

		if tokenType == "refresh" {
			secret = jwtRefreshSecret
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("could not parse token")
	}

	isTokenValid := parsedToken.Valid

	if !isTokenValid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userId, ok := claims["userId"].(float64)

	if !ok {
		return nil, errors.New("invalid userId in token claims")
	}

	email, ok := claims["email"].(string)

	if !ok {
		return nil, errors.New("invalid email in token claims")
	}

	return &TokenClaims{
		UserId: int64(userId),
		Email:  email,
	}, nil
}
