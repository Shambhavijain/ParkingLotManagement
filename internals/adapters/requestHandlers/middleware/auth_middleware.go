package middleware

import (
	"net/http"
	"parkingSlotManagement/internals/core/services/auth"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc, authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Split "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		// fmt.Println("Clean token received:", tokenStr)

		_, err := authService.ValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
