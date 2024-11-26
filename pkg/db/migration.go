package db

import (
	"database/sql"
	"log"
)

// RunMigrations applies the required database schema migrations.
func RunMigrations(DB *sql.DB) {
	usersTable := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            email_verified BOOLEAN DEFAULT FALSE,
            verification_token VARCHAR(255),
            role ENUM('admin', 'user') DEFAULT 'user',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatalf("Migration failed for 'users' table: %v", err)
	}

	activityLogsTable := `
        CREATE TABLE IF NOT EXISTS activity_logs (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            action VARCHAR(255) NOT NULL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `
	_, err = DB.Exec(activityLogsTable)
	if err != nil {
		log.Fatalf("Migration failed for 'activity_logs' table: %v", err)
	}

	log.Println("Migration completed successfully.")
}
