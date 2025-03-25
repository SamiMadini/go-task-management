package main

import (
	commons "sama/go-task-management/commons"
	"sama/go-task-management/gateway/middleware"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_EXPIRATION = 10
const REFRESH_TOKEN_EXPIRATION = 30

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func generateTokens(userID string) (*TokenResponse, error) {
	accessToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, &middleware.AuthClaims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(ACCESS_TOKEN_EXPIRATION) * time.Minute)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	})

	refreshToken := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, &middleware.AuthClaims{
		UserID: userID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRATION) * 24 * time.Hour)), // 30 days
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
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
