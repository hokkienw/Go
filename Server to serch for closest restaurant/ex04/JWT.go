package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("1234")

func main() {
	http.HandleFunc("/api/get_token", GetToken)
	http.HandleFunc("/api/recommend", ProtectedEndpoint)

	http.ListenAndServe(":8888", nil)
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := generateToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"token": "%s"}`, tokenString)
}

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	token, err := validateToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Protected API endpoint")
}

func generateToken() (string, error) {
	claims := jwt.MapClaims{
		"name":  "hokkienw",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(r *http.Request) (*jwt.Token, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("Authorization header missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
