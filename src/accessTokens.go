package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserId string `json:"userId"` // hex representation of mongodb ObjectID
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var hmacSecret = []byte(os.Getenv("JWT_SECRET"))

func GetAccessTokenForSession(email string, userId string) (string, error) {
	expTime := time.Now().Add(1 * time.Hour)
	claims := &AccessTokenClaims{
		Email:  email,
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func ValidateAccessToken(tokenString string, c *AccessTokenClaims) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Email = claims["email"].(string)
		c.UserId = claims["userId"].(string)
		return nil
	}
	return fmt.Errorf("unable to parse the token")
}
