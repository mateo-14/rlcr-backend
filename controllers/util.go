package controllers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

/* func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	hash, _ := hashPassword("pass")

	lib.FirestoreClient.Collection("settings").Doc("user").Set(r.Context(), map[string]interface{}{
		"user":     "admin",
		"password": hash,
	})

	fmt.Fprintf(w, "Account updated")
} */

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(expireIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ExpiresAt: time.Now().Add(expireIn).Unix()})
	tsecret := os.Getenv("TOKEN_SECRET")

	return token.SignedString([]byte(tsecret))
}
