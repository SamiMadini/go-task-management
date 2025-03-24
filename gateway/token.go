package main

import (
	commons "sama/go-task-management/commons"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_EXPIRATION = 10
const REFRESH_TOKEN_EXPIRATION = 30

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func generateTokens(userID string) (*TokenResponse, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ACCESS_TOKEN_EXPIRATION) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRATION) * 24 * time.Hour)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	secretKey := []byte(commons.GetEnv("JWT_SECRET", "your-secret-key"))
	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    ACCESS_TOKEN_EXPIRATION * 60,
	}, nil
}
