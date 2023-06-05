package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

type auth interface {
	VerifyToken(token string) (userID int, err error)
	UserPermissions(id int) (string, error)
}

type accessChecker struct {
	a auth
}

func (a accessChecker) checkAccessRights(next http.Handler, perms string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "invalid authorization header", http.StatusBadRequest)
			return
		}

		userID, err := a.a.VerifyToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintln("could not authorize:", err), http.StatusUnauthorized)
			return
		}

		userPermissions, err := a.a.UserPermissions(userID)
		if err != nil {
			http.Error(w, "failed to get user permissions", http.StatusInternalServerError)
			return
		}

		if perms != userPermissions {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
