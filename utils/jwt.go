package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signMethodErrorMessage = "Unexpected signing method"
var invalidTokenErrorMessage = "Invalid token"
var invalidAuthHeaderErrorMessage = "Invalid auth header"

func GenerateJwtToken(userId int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func VerifyJwtToken(authHeader string) (int64, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, errors.New(invalidAuthHeaderErrorMessage)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New(signMethodErrorMessage)
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, errors.New(invalidTokenErrorMessage)
	}

	if !parsedToken.Valid {
		return 0, errors.New(invalidTokenErrorMessage)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New(invalidTokenErrorMessage)
	}

	userId, ok := claims["userId"].(float64)

	if !ok {
		return 0, errors.New(invalidTokenErrorMessage)
	}

	return int64(userId), nil
}
