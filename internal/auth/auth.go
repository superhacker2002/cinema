package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

// User - модель данных для пользователей
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Token - модель данных для JWT-токена
type Token struct {
	Token string `json:"token"`
}

// JWTSecret - секретный ключ для создания и проверки JWT-токена
var JWTSecret = []byte("my-secret-key")

// loginHandler - обработчик для авторизации пользователя и создания JWT-токена
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ищем пользователя в базе данных
	//row := DB.QueryRow("SELECT id, password FROM users WHERE username=$1", user.Username)
	//var dbUser User
	//err = row.Scan(&dbUser.ID, &dbUser.Password)
	//if err != nil {
	//	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	//	return
	//}

	// Проверяем пароль пользователя
	//if user.Password != dbUser.Password {
	//	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	//	return
	//}

	// Создаем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем JWT-токен
	response := Token{Token: tokenString}
	json.NewEncoder(w).Encode(response)
}

//
//// getUsersHandler - обработчик для получения списка пользователей
//func getUsersHandler(w http.ResponseWriter, r *http.Request) {
//	// Получаем токен из заголовка Authorization
//	authHeader := r.Header.Get("Authorization")
//	if authHeader == "" {
//		http.Error(w, "Authorization header is required", http.StatusBadRequest)
//		return
