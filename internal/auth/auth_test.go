package auth

import (
	userRepository "bitbucket.org/Ernst_Dzeravianka/cinemago-app/internal/user/repository"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockRepository struct {
	creds userRepository.Credentials
	err   error
}

func (m mockRepository) User(username string) (userRepository.Credentials, error) {
	return m.creds, m.err
}

func TestAuthenticate(t *testing.T) {
	repo := mockRepository{}

	hasher := sha256.New()
	hasher.Write([]byte("password"))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	repo.creds = userRepository.Credentials{"existing_user", hashedPassword}

	t.Run("Valid auth", func(t *testing.T) {
		repo.err = nil
		auth := New("secret-key", 24, repo)
		token, err := auth.Authenticate("existing_user", hashedPassword)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid username", func(t *testing.T) {
		repo.err = ErrInvalidUsernameOrPassword
		auth := New("secret-key", 24, repo)
		token, err := auth.Authenticate("non_existing_user", hashedPassword)
		assert.Equal(t, ErrInvalidUsernameOrPassword, err)
		assert.Empty(t, token)
	})

	t.Run("Invalid password", func(t *testing.T) {
		repo.err = ErrInvalidUsernameOrPassword
		auth := New("secret-key", 24, repo)
		token, err := auth.Authenticate("existing_user", "invalid_password")
		assert.Equal(t, ErrInvalidUsernameOrPassword, err)
		assert.Empty(t, token)
	})
}

func createTokenString(secret []byte, userID string, tokenExp int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(tokenExp)).Unix(),
	})

	return token.SignedString(secret)
}

func TestVerifyToken(t *testing.T) {
	repo := mockRepository{}

	repo.creds = userRepository.Credentials{}
	auth := New("secret-key", 24, repo)

	t.Run("valid token", func(t *testing.T) {
		token, _ := createTokenString([]byte("secret-key"), "existing_user", 24)
		userID, err := auth.VerifyToken(token)

		assert.Nil(t, err, "unexpected error occurred: %w", err)
		assert.Equal(t, "existing_user", userID)
	})

	t.Run("invalid token", func(t *testing.T) {
		invalidToken := "invalid_token"
		userID, err := auth.VerifyToken(invalidToken)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Equal(t, "", userID,
			"user id should be empty when token is invalid")
	})

	t.Run("invalid signing method", func(t *testing.T) {
		invalidToken, _ := createTokenString([]byte("invalid-secret-key"), "existing_user", 24)
		userID, err := auth.VerifyToken(invalidToken)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Equal(t, "", userID,
			"user id should be empty when token was signed by invalid method")
	})

	t.Run("expired token", func(t *testing.T) {
		auth.exp = 0
		token, _ := createTokenString([]byte("secret-key"), "existing_user", 0)
		userID, err := auth.VerifyToken(token)

		assert.Equal(t, ErrExpiredToken, err)
		assert.Empty(t, userID, "user id should be empty when token is expired")
	})
}
