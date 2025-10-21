package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Read secret from environment with a safe fallback
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	// Build claims - adjust as needed
	claims := jwt.MapClaims{
		"sub":  "stocky",
		"role": "user",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iss":  "stocky",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sign token: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Authorization: Bearer %s\n", signedToken)
}
