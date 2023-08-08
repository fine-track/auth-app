package utils

import (
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

// TODO: add OTP verification methods
