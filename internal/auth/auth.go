package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrInvalidSigningMethod      = errors.New("invalid signing method")
	ErrInvalidToken              = errors.New("invalid token")
)

type Credentials struct {
	ID           string
	PasswordHash string
}

type repository interface {
	User(username string) (Credentials, error)
}

type auth struct {
	jwtSecret []byte
	r         repository
}

func New(jwtSecret string, repository repository) auth {
	return auth{
		jwtSecret: []byte(jwtSecret),
		r:         repository,
	}
}

func (a auth) Authenticate(username string, passwordHash string) (string, error) {
	userCreds, err := a.r.User(username)
	if err != nil {
		return "", err
	}

	if passwordHash != userCreds.PasswordHash {
		return "", ErrInvalidUsernameOrPassword
	}

	return a.generateJWT(userCreds.ID)
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

func (a auth) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		log.Println(err)
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", ErrInvalidToken
}
