package authmw

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type auth interface {
	VerifyToken(token string) (userID int, err error)
	UserPermissions(id int) (string, error)
}

type AccessChecker struct {
	a auth
}

func New(a auth) AccessChecker {
	return AccessChecker{
		a: a,
	}
}

func (a AccessChecker) Authenticate(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a AccessChecker) CheckPerms(next http.Handler, perms ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		userPermissions, err := a.a.UserPermissions(userID)
		if err != nil {
			http.Error(w, "failed to get user permissions", http.StatusInternalServerError)
			return
		}

		if !hasPermissions(userPermissions, perms) {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hasPermissions(userPerms string, reqPerms []string) bool {
	for _, word := range reqPerms {
		if word == userPerms {
			return true
		}
	}
	return false
}
