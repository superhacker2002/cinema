package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type Credentials struct {
	ID           string
	PasswordHash string
}

type repository interface {
	GetUserInfo(username string) (Credentials, error)
}

type auth struct {
	jwtSecret  []byte
	repository repository
}

func New(jwtSecret string, repository repository) auth {
	return auth{
		jwtSecret:  []byte(jwtSecret),
		repository: repository,
	}
}

func (a auth) Authenticate(username string, password string) (string, error) {
	userInfo, err := a.repository.GetUserInfo(username)
	if err != nil {
		return "", err
	}

	err = a.comparePasswords(userInfo.PasswordHash, []byte(password))
	if err != nil {
		return "", err
	}

	return a.generateJWT(userInfo.ID)
}

func (a auth) generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a auth) VerifyToken(tokenString string) error {
	// TODO: verify and decode token
	return nil
}

func (a auth) comparePasswords(hash string, password []byte) error {
	hasher := sha256.New()
	hasher.Write(password)
	passwordHash := hex.EncodeToString(hasher.Sum(nil))
	if hash != passwordHash {
		return ErrInvalidPassword
	}
	return nil
}
