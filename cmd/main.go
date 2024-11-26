package main

import (
	"log"
	"net/http"
	"user-management-system/api"
	"user-management-system/config"
	"user-management-system/pkg/db"
)

func main() {
	config.LoadEnv()
	database := db.InitDB()    // Assign the return value of db.InitDB() to a variable
	db.RunMigrations(database) // Pass the database instance to RunMigrations

	router := api.SetupRoutes()
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
