package database

import (
	"database/sql"
	"log"
)

func Migrate(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255),
			email VARCHAR(255) UNIQUE,
			password VARCHAR(255),
			role VARCHAR(50),
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS services (
			id INT AUTO_INCREMENT PRIMARY KEY,
			staff_id INT,
			name VARCHAR(255),
			duration INT,
			price DECIMAL(10,2),
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS appointments (
			id INT AUTO_INCREMENT PRIMARY KEY,
			client_id INT,
			staff_id INT,
			service_id INT,
			scheduled_at DATETIME,
			status VARCHAR(50),
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS available_slots (
			id INT AUTO_INCREMENT PRIMARY KEY,
			staff_id INT,
			weekday VARCHAR(20),
			start_time TIME,
			end_time TIME
		)`,
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			log.Fatalf("erro ao executar migration: %v", err)
		}
	}
}
