package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID uint, role string) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id" : userID,
		"role": role,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(JWT_SECRET)
}

func ValidateToken(tokenString string) (*jwt.Token, error){
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return JWT_SECRET, nil
	})
}