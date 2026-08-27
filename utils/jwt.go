package utils

import (
	"crypto/sha256"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AccessSecret))
}

func GenerateRefreshToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshSecret))
}

func ValidateAccessToken(tokenStr string) (string, error, bool) {
	// Returns (username, err, isExpired)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(AccessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", err, true
		}
		// Also check error string for expiration or time is in the past
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return "", err, false
		}
		// Sometimes jwt package returns err directly if expired
		return "", err, true // or check if expired
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil, false
		}
	}
	return "", errors.New("invalid token"), false
}

func ValidateRefreshToken(tokenStr string) (string, error, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(RefreshSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", err, true
		}
		return "", err, true
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil, false
		}
	}
	return "", errors.New("invalid refresh token"), false
}

func GetAESKey() []byte {
	hash := sha256.Sum256([]byte(ConfigSecret))
	return hash[:]
}
