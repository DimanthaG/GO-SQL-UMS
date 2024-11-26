package handlers

import (
	"database/sql"
)

var DB *sql.DB

// InitHandlers initializes the handlers with the provided database instance.
func InitHandlers(db *sql.DB) {
	DB = db
}
