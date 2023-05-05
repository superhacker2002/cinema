package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepository struct {
	creds Credentials
	err   error
}

func (m mockRepository) GetUserInfo(username string) (Credentials, error) {
	return m.creds, m.err
}

func TestAuthenticate(t *testing.T) {
	repo := mockRepository{}

	hasher := sha256.New()
	hasher.Write([]byte("password"))
	hashed_password := hex.EncodeToString(hasher.Sum(nil))

	repo.creds = Credentials{"existing_user",
		hashed_password}

	t.Run("Valid auth", func(t *testing.T) {
		repo.err = nil
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("existing_user", "password")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid username", func(t *testing.T) {
		repo.err = ErrUserNotFound
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("non_existing_user", "password")
		assert.Equal(t, ErrUserNotFound, err)
		assert.Empty(t, token)
	})

	t.Run("Invalid password", func(t *testing.T) {
		repo.err = ErrInvalidPassword
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("existing_user", "invalid_password")
		assert.Equal(t, ErrInvalidPassword, err)
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
		assert.Equal(t, ErrInvalidPassword, err)
	})
}
