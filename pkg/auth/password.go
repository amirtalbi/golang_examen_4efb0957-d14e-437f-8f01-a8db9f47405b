package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    log.Printf("Checking password hash for user: %s", password)
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    log.Printf("Password check result: %v", err == nil)
    return err == nil
}