package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "file:db.sqlite")

	if err != nil {
		log.Fatalf("Failed to establish connection to database: %v", err)
	}

	DB = db

	// exec create tables here
	generateTableFromStatement(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		fullname TEXT DEFAULT ''
	)`)

	return db
}

// create tables here
func generateTableFromStatement(stmt string) {
	_, err := DB.Exec(stmt)

	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}
}
