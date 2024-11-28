package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"user-management-system/pkg/auth"
	"user-management-system/pkg/db"
)

// SignIn Handler
func SignIn(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the JSON input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	// Variables for database query
	var hashedPassword string
	var emailVerified bool

	// Query user data from the database
	err := db.DB.QueryRow("SELECT password_hash, email_verified FROM users WHERE email = ?", req.Email).
		Scan(&hashedPassword, &emailVerified)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Printf("Error querying user data: %v", err)
		return
	}

	// Check if email is verified
	if !emailVerified {
		http.Error(w, "Email not verified", http.StatusUnauthorized)
		log.Printf("Email not verified for user: %s", req.Email)
		return
	}
	// Check the password
	if err := auth.CheckPassword(hashedPassword, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Printf("Password mismatch for user: %s", req.Email)
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateJWT(req.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Error generating JWT: %v", err)
		return
	}

	// Send the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
