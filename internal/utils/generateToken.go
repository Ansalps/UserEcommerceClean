package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email string, ID uint, role string, expiry uint) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"ID":    ID,
		"exp":   time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		"role":  role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
