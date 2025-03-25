package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwtv5.RegisteredClaims
}

type AuthConfig struct {
	JWTSecret     string
	PublicPaths   []string
	SwaggerPrefix string
}

func DefaultAuthConfig(jwtSecret string) AuthConfig {
	return AuthConfig{
		JWTSecret: jwtSecret,
		PublicPaths: []string{
			"/api/_health",
			"/api/v1/auth/signin",
			"/api/v1/auth/signup",
			"/api/v1/auth/refresh-token",
			"/api/v1/auth/forgot-password",
			"/api/v1/auth/reset-password",
		},
		SwaggerPrefix: "/swagger/",
	}
}

func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			if isPublicPath(r.URL.Path, config) {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := validateAuthHeader(r, config)
			if err != nil {
				fmt.Printf("AuthMiddleware: Authentication failed: %v\n", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicPath(path string, config AuthConfig) bool {
	if strings.HasPrefix(path, config.SwaggerPrefix) {
		return true
	}

	path = strings.TrimSuffix(path, "/")
	for _, publicPath := range config.PublicPaths {
		publicPath = strings.TrimSuffix(publicPath, "/")
		if path == publicPath {
			return true
		}
	}

	return false
}

func validateAuthHeader(r *http.Request, config AuthConfig) (*AuthClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("JWT::Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("JWT::Invalid authorization header format")
	}

	token, err := ParseToken(parts[1], config.JWTSecret)
	if err != nil {
		fmt.Printf("validateAuthHeader: Token validation failed: %v\n", err)
		return nil, fmt.Errorf("JWT::Invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		fmt.Printf("validateAuthHeader: Invalid token claims\n")
		return nil, fmt.Errorf("JWT::Invalid token claims")
	}

	return claims, nil
}

func ParseToken(tokenString, secret string) (*jwtv5.Token, error) {
	return jwtv5.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwtv5.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
