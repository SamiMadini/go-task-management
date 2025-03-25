package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwtv5.RegisteredClaims
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

func parseToken(tokenString string) (*jwtv5.Token, error) {
	return jwtv5.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}
