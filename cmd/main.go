package main

import (
	"log"
	"net/http"
	"user-management-system/api"
	"user-management-system/config"
	"user-management-system/pkg/db"
	"user-management-system/pkg/handlers"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	database := db.InitDB()

	// Run database migrations
	db.RunMigrations(database)

	// Pass the database instance to handlers
	handlers.InitHandlers(database)

	// Set up routes and start server
	router := api.SetupRoutes()
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
