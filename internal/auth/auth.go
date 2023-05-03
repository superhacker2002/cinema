package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Auth struct {
	db        *sql.DB
	jwtSecret []byte
}

func New(db *sql.DB, jwtSecret string) Auth {
	return Auth{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

func (a Auth) Authenticate(username string, password string) (string, error) {
	var passwordHash string
	var userID int64

	// TODO: get user id from db

	err := comparePasswords(passwordHash, []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return a.generateJWT(userID)
}

func (a Auth) generateJWT(userID int64) (string, error) {
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

func (a Auth) VerifyToken(tokenString string) (*User, error) {
	// TODO: verify and decode token
	return &User{}, nil
}

func comparePasswords(hash string, password []byte) error {
	// TODO: compare password with its hash
	return nil
}
