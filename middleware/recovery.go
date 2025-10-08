package middleware

import (
	"encoding/json"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "SERVER_ERROR",
					"message": "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
