package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() *sql.DB {
	dsn := os.Getenv("MYSQL_DSN")
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	log.Println("Database connection established")
	return DB
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}
}
