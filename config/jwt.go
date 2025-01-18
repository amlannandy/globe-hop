package config

import (
	"globe-hop/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(user *models.User) (string, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")
	userId := user.ID
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token will be valid for 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

func DecodeJWTToken(tokenString string) (float64, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")

	// Parse the token
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure that the signing mathed is as expected
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecretKey), nil
	})

	// Handle parsing errors
	if err != nil {
		return -1, err
	}

	// Validate the claims and extract them.
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if ok && parsedToken.Valid {
		// Extract user_id from claims
		userId, ok := claims["user_id"]

		if !ok {
			return -1, jwt.ErrInvalidKey
		}

		return userId.(float64), nil
	}

	return -1, jwt.ErrInvalidKey
}
