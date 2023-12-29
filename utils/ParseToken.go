package utils

import (
	"os"
	"server/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func ParseToken(tokenString string) (claims *models.Claims, erro error) {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
