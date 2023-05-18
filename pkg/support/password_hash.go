package support

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// password hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) //GenerateFromPassword returns the bcrypt hash of the password
	return string(bytes), err
}

// password authorization
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) //CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
	fmt.Println("password", err, password, hash)
	if err == nil {
		// return true
		return nil
	} else {
		// return false
		return errors.New("invalid password")
	}
}
