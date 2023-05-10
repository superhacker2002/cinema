package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepository struct {
	creds Credentials
	err   error
}

func (m mockRepository) User(username string) (Credentials, error) {
	return m.creds, m.err
}

func TestAuthenticate(t *testing.T) {
	repo := mockRepository{}

	hasher := sha256.New()
	hasher.Write([]byte("password"))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	repo.creds = Credentials{"existing_user", hashedPassword}

	t.Run("Valid auth", func(t *testing.T) {
		repo.err = nil
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("existing_user", "password")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid username", func(t *testing.T) {
		repo.err = ErrInvalidUsernameOrPassword
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("non_existing_user", "password")
		assert.Equal(t, ErrInvalidUsernameOrPassword, err)
		assert.Empty(t, token)
	})

	t.Run("Invalid password", func(t *testing.T) {
		repo.err = ErrInvalidUsernameOrPassword
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("existing_user", "invalid_password")
		assert.Equal(t, ErrInvalidUsernameOrPassword, err)
		assert.Empty(t, token)
	})
}

func TestComparePasswords(t *testing.T) {
	auth := auth{}
	password := "password"
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hash := hex.EncodeToString(hasher.Sum(nil))

	t.Run("Compare valid password", func(t *testing.T) {
		err := auth.comparePasswords(hash, []byte(password))
		assert.Nil(t, err)
	})

	t.Run("Compare invalid password", func(t *testing.T) {
		err := auth.comparePasswords(hash, []byte("invalid_password"))
		assert.Equal(t, ErrInvalidUsernameOrPassword, err)
	})
}

func createTokenString(secret []byte, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	return token.SignedString(secret)
}

func TestVerifyToken(t *testing.T) {
	repo := mockRepository{}

	repo.creds = Credentials{}
	auth := New("secret-key", repo)

	t.Run("Valid token", func(t *testing.T) {
		token, _ := createTokenString([]byte("secret-key"), "existing_user")
		userID, err := auth.VerifyToken(token)

		assert.Nil(t, err, "unexpected error occurred: %w", err)
		assert.Equal(t, "existing_user", userID)
	})

	t.Run("Invalid token", func(t *testing.T) {
		invalidToken := "invalid_token"
		userID, err := auth.VerifyToken(invalidToken)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Equal(t, "", userID)
	})

	t.Run("Invalid signing method", func(t *testing.T) {
		invalidToken, _ := createTokenString([]byte("invalid-secret-key"), "existing_user")
		userID, err := auth.VerifyToken(invalidToken)

		assert.Equal(t, ErrInvalidToken, err)
		assert.Equal(t, "", userID)
	})
}
