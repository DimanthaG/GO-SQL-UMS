package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := sql.Open("mysql", dsn) // Change variable to lowercase
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	log.Println("Database connection established")
	return db // Return the database instance
}
