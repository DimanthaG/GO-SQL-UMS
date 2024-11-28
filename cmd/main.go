package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtKey = []byte("your_secret_key") // Replace with a strong secret key

// Initialize the database connection
func initDB() {
	var err error
	dsn := "app_user:app_password@tcp(db:3306)/user_management"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	log.Println("Database connection established")
}

// Validate email format using regex
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// Signup handler
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if len(user.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	_, err = db.Exec(`
        INSERT INTO users (email, password_hash, email_verified, role)
        VALUES (?, ?, ?, ?)`,
		user.Email, hashedPassword, 0, "user")
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) && mySQLErr.Number == 1062 { // Duplicate entry
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		log.Printf("Database error during signup: %v\n", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	log.Printf("User %s signed up successfully\n", user.Email)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Signup successful")
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = ?", creds.Email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			log.Printf("Database error during login: %v\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	expirationTime := time.Now().Add(30 * time.Minute) // Token expires in 30 minutes
	claims := &jwt.MapClaims{
		"email": creds.Email,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Set token as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true, // Prevents access to the cookie via JavaScript
	})

	log.Printf("User %s logged in successfully\n", creds.Email)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login successful")
}
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the token cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Parse and validate the token
		tokenStr := cookie.Value
		claims := &jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}

// GenerateRandomPassword generates a random password of a given length
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	for i := range password {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[random.Int64()]
	}
	return string(password), nil
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the token cookie by setting it to an empty value with an expired date
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expired date
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Logout successful")
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email string `json:"email"`
	}

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Email == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Invalid input: %v", err)
		return
	}

	// Check if the email exists in the database
	var email string
	err = db.QueryRow("SELECT email FROM users WHERE email = ?", request.Email).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Email not found", http.StatusNotFound)
			log.Printf("Email not found: %s", request.Email)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Printf("Database error: %v", err)
		}
		return
	}

	log.Printf("Email verified: %s", email)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Email verified successfully")
}

// Reset Password Handler
func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email           string `json:"email"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	// Decode and validate the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request.Email == "" || request.CurrentPassword == "" || request.NewPassword == "" {
		log.Println("Invalid input: ", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = ?", request.Email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Email not found: %s\n", request.Email)
			http.Error(w, "Email not found", http.StatusNotFound)
		} else {
			log.Printf("Database error for email %s: %v\n", request.Email, err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Validate the current password
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(request.CurrentPassword))
	if err != nil {
		log.Printf("Invalid current password for email: %s\n", request.Email)
		http.Error(w, "Invalid current password", http.StatusUnauthorized)
		return
	}

	// Check the strength of the new password
	if len(request.NewPassword) < 8 {
		log.Printf("Weak password provided for email: %s\n", request.Email)
		http.Error(w, "New password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing new password for email %s: %v\n", request.Email, err)
		http.Error(w, "Error hashing new password", http.StatusInternalServerError)
		return
	}

	// Update the password in the database
	_, err = db.Exec("UPDATE users SET password_hash = ? WHERE email = ?", hashedPassword, request.Email)
	if err != nil {
		log.Printf("Failed to update password for email %s: %v\n", request.Email, err)
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	log.Printf("Password updated successfully for email: %s\n", request.Email)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password updated successfully")
}
func videoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the video page!")
}

// Enable CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust "*" for security in production
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	initDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/verify-email", verifyEmailHandler)
	mux.HandleFunc("/reset-password", resetPasswordHandler)
	mux.Handle("/video", authenticate(http.HandlerFunc(videoHandler)))
	mux.HandleFunc("/logout", logoutHandler)

	server := &http.Server{Addr: ":8080", Handler: enableCORS(mux)}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		server.Close()
		db.Close()
		os.Exit(0)
	}()

	log.Println("Server is running on port 8080")
	log.Fatal(server.ListenAndServe())
}
