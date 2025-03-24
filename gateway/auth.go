package main

import (
	"context"
	"net/http"
	"strings"

	commons "sama/go-task-management/commons"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const userIDKey contextKey = "user_id"

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/_health" ||
		   strings.HasPrefix(r.URL.Path, "/swagger/") ||
		   strings.HasPrefix(r.URL.Path, "/api/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "JWT::Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "JWT::Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := parseToken(parts[1])
		if err != nil {
			http.Error(w, "JWT::Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*AuthClaims)
		if !ok || !token.Valid {
			http.Error(w, "JWT::Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(commons.GetEnv("JWT_SECRET", "your-secret-key"))

	return jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}
