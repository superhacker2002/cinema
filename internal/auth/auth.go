package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

type Credentials struct {
	ID       string
	Password string
}

type repository interface {
	getUserInfo(username string) (Credentials, error)
}

type auth struct {
	jwtSecret  []byte
	repository repository
}

func New(jwtSecret string) auth {
	return auth{
		jwtSecret: []byte(jwtSecret),
	}
}

func (a auth) Authenticate(username string, password string) (string, error) {
	userInfo, err := a.repository.getUserInfo(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = a.comparePasswords(userInfo.Password, []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
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
	// TODO: compare password with its hash
	return nil
}
