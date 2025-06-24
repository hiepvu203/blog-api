package utils

import (
	// "errors"
	"os"
	// "strconv"
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

// func GenerateResetToken(userID uint, duration time.Duration) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": userID,
// 		"type": "reset",
// 		"exp": time.Now().Add(duration).Unix(),
// 	})
// 	return token.SignedString(JWT_SECRET)
// }

// func ValidateResetToken(tokenString string) (uint, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrTokenSignatureInvalid
// 		}
// 		return JWT_SECRET, nil
// 	})
// 	if err != nil {
// 		return 0, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok || !token.Valid {
// 		return 0, errors.New("invalid token")
// 	}

// 	if t, ok := claims["type"].(string); !ok || t != "reset" {
// 		return 0, errors.New("invalid token type")
// 	}

// 	var userID uint
// 	switch v := claims["user_id"].(type) {
// 	case float64:
// 		userID = uint(v)
// 	case string:
// 		parsed, err := strconv.ParseUint(v, 10, 64)
// 		if err != nil {
// 			return 0, errors.New("invalid user_id in token")
// 		}
// 		userID = uint(parsed)
// 	default:
// 		return 0, errors.New("invalid user_id in token")
// 	}

// 	return userID, nil
// }