package main

import (
	"log"
	"net/http"
	"user-management-system/api"
	"user-management-system/pkg/db"
	"user-management-system/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db.InitDB()

	// Run database migrations
	db.RunMigrations()

	// Setup routes
	router := api.SetupRoutes()

	// Start the server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
