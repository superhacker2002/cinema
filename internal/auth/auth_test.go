package auth

import (
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

func TestAuth_Authenticate(t *testing.T) {
	repo := mockRepository{}

	t.Run("Valid auth", func(t *testing.T) {
		repo.creds = Credentials{"existing_user", "password"}
		repo.err = nil
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("existing_user", "password")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid username", func(t *testing.T) {
		repo.creds = Credentials{"existing_user", "password"}
		repo.err = ErrUserNotFound
		auth := New("secret-key", repo)
		token, err := auth.Authenticate("non_existing_user", "password")
		assert.Equal(t, ErrUserNotFound, err)
		assert.Empty(t, token)
	})

	//// Test invalid password
	//token, err = auth.Authenticate("existing_user", "wrong_password")
	//assert.Equal(t, errors.New("invalid password"), err, "Error should be 'invalid password'")
	//assert.Empty(t, token, "Token should be empty")
}
