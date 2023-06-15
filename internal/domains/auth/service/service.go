package service

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
	ErrExpiredToken              = errors.New("token is expired")
	ErrInternalError             = errors.New("internal server error")
	ErrUserNotFound              = errors.New("user not found")
)

type Credentials struct {
	ID           int
	PasswordHash string
}

type repository interface {
	GetUser(username string) (Credentials, error)
	Permissions(userId int) (string, error)
}

type Auth struct {
	jwtSecret []byte
	r         repository
	exp       int
}

func New(jwtSecret string, tokenExp int, repo repository) Auth {
	return Auth{
		jwtSecret: []byte(jwtSecret),
		r:         repo,
		exp:       tokenExp,
	}
}

func (a Auth) Authenticate(username string, passwordHash string) (token string, err error) {
	userCreds, err := a.r.GetUser(username)
	if errors.Is(ErrUserNotFound, err) {
		log.Println(err)
		return "", ErrInvalidUsernameOrPassword
	}
	if err != nil {
		log.Println(err)
		return "", ErrInternalError
	}

	if passwordHash != userCreds.PasswordHash {
		log.Println(err)
		return "", ErrInvalidUsernameOrPassword
	}

	return a.generateJWT(userCreds.ID)
}

func (a Auth) generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(a.exp)).Unix(),
	})
	signedToken, err := token.SignedString(a.jwtSecret)
	if err != nil {
		log.Println("failed to sign token:", err)
		return "", err
	}
	return signedToken, nil
}

func (a Auth) VerifyToken(token string) (userID int, err error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return a.jwtSecret, nil
	})
	if err != nil {
		log.Println("failed to parse token:", err)
		return 0, ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if ok && parsedToken.Valid {
		if a.tokenIsExpired(claims) {
			return 0, ErrExpiredToken
		}
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, ErrInvalidToken
}

func (a Auth) tokenIsExpired(claims jwt.MapClaims) bool {
	exp := time.Unix(int64(claims["exp"].(float64)), 0).UTC()
	now := time.Now().UTC()
	return exp.Before(now)
}

func (a Auth) UserPermissions(id int) (string, error) {
	perms, err := a.r.Permissions(id)
	if errors.Is(err, ErrUserNotFound) {
		log.Println("failed to get user permissions:", ErrUserNotFound)
		return "", err
	}
	if err != nil {
		log.Println("failed to get user permissions:", err)
		return "", ErrInternalError
	}

	return perms, nil
}
