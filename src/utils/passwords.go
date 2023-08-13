package utils

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// return the hash for the password
func HashPass(p string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// the first arg is the `password` and the second one is the `hash` string
func IsValidPass(p string, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}

func GenRandomStr(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // using all uppercase for better user experience
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
