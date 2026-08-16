package jwtauth

import (
	"ab_system/pkg/errs"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ИСПОЛЬЗОВАНО СО СТАРОГО ПРОЕКТА

type CustomClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(userID string,
	role string,
	appName string,
	duration time.Duration) (string, error) {

	claims := &CustomClaims{
		Sub:  userID,
		Role: role,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    appName,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("RANDOM_SECRET")))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("RANDOM_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errs.ErrInvalidToken
}
