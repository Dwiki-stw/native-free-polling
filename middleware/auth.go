package middleware

import (
	"context"
	"encoding/json"
	"native-free-pollings/helper"
	"net/http"
	"strings"
	"time"
)

type ctxKey string

const (
	userIDKey ctxKey = "userID"
)

func Auth(screet []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "INVALID_TOKEN",
					"message": "missing or invalid token",
				})
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := helper.ExtractToken(token, screet)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "INVALID_TOKEN",
					"message": err.Error(),
				})
				return
			}

			now := time.Now()
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "EXPIRED_TOKEN",
					"message": "token expired",
				})
				return
			}

			if claims.NotBefore != nil && claims.NotBefore.After(now) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "NOT_VALID",
					"message": "token not yet valid",
				})
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
