package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

func GenerateJWT(JWTSecret []byte) (Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return Token{}, err
	}

	return Token{Token: tokenString}, nil
}
