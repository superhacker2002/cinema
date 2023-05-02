package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	token, err := generateJWT()
	assert.NoError(t, err)

	parsedToken, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	assert.NoError(t, err)

	assert.True(t, parsedToken.Valid)
}

func TestJWTSecret(t *testing.T) {
	assert.NotEmpty(t, JWTSecret)
	assert.NotNil(t, JWTSecret)
}
