package main

import (
	"log"
	"net/http"
	"os"
	"user-management-system/api"
	"user-management-system/config"
	"user-management-system/pkg/db"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Set up routes
	router := api.SetupRoutes()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
