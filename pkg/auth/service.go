package auth

import (
	"database/sql"
	"errors"

	"github.com/sina-byn/go-jwt-auth/pkg/user"
	"github.com/sina-byn/go-jwt-auth/pkg/utils"
)

func Login(email, password string) (*utils.TokenPair, error) {
	user, err := user.GetUserByEmail(email)

	if user == nil && err == nil {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	err = utils.ValidateHash(password, user.Password)

	if err != nil {
		return nil, errors.New("invalid password")
	}

	tokenPair, err := utils.GenerateTokenPair(user.Id, email)

	return tokenPair, err
}

func Refresh(refreshToken string) (*utils.TokenPair, error) {
	tokenClaims, err := utils.VerifyToken(refreshToken, "refresh")

	if err != nil {
		return nil, err
	}

	refreshedTokenPair, err := utils.GenerateTokenPair(tokenClaims.UserId, tokenClaims.Email)

	return refreshedTokenPair, err
}
