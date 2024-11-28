package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateRandomToken creates a secure random token
func GenerateRandomToken() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}
	return hex.EncodeToString(bytes)
}
