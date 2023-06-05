package middleware

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header required")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
)

type auth interface {
	VerifyToken(token string) (userID int, err error)
}

func checkAccessRights(next http.Handler, perms []string, a auth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusBadRequest)
			return
		}

		//userID, err := a.VerifyToken(token)
		//if err != nil {
		//	http.Error(w, fmt.Sprintln("could not authorize:", err), http.StatusUnauthorized)
		//	return
		//}

		// Получение прав пользователя из базы данных или другого источника
		//userPermissions, err := getUserPermissions(userID)
		//if err != nil {
		//	http.Error(w, "Failed to get user permissions", http.StatusInternalServerError)
		//	return
		//}
		//
		//// Проверка прав пользователя для доступа к эндпоинту
		//if !hasRequiredPermissions(userPermissions, requiredPermissions) {
		//	http.Error(w, "Insufficient permissions", http.StatusForbidden)
		//	return
		//}

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)
	})
}
