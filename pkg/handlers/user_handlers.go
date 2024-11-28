package handlers

import (
	"encoding/json"
	"net/http"
	"user-management-system/pkg/auth"
	"user-management-system/pkg/db"
	"user-management-system/pkg/email"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	verificationToken := auth.GenerateRandomToken()

	_, err = db.DB.Exec("INSERT INTO users (email, password_hash, verification_token) VALUES (?, ?, ?)",
		req.Email, hashedPassword, verificationToken)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	go email.SendEmail(req.Email, "Verify Your Email",
		"Click the link to verify your email: http://localhost:8080/verify-email?token="+verificationToken)

	json.NewEncoder(w).Encode(map[string]string{"message": "User registered. Please verify your email."})
}
