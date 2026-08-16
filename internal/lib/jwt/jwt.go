package jwtauth

import (
	"os"
	"strconv"
	"time"
)

func NewToken(userID string,
	role string,
	appName string) (string, error) {

	jwtTtlSecondsInt, err := strconv.ParseInt(os.Getenv("JWT_TTL_SECONDS"), 10, 64)
	if err != nil {
		return "", err
	}

	accessToken, err := GenerateToken(userID, role, appName, time.Duration(jwtTtlSecondsInt)*time.Second)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
