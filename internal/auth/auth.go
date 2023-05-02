package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

// JWTSecret should be generated randomly and saved in environment var
var JWTSecret = []byte("my-secret-key")

func generateJWT() (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return Token{}, err
	}

	return Token{Token: tokenString}, nil
}
