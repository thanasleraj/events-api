package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "events-api.db")

	if err != nil {
		panic("Failed to establish database connection")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	enableForeignKeys := "PRAGMA foreign_keys = ON"

	_, err := DB.Exec(enableForeignKeys)
	if err != nil {
		panic("Failed to enable foreign keys")
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		panic("Failed to create users table")
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		userId INTEGER,
		FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
	)
	`
	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic("Failed to create events table")
	}

	createIndex := `
	CREATE INDEX IF NOT EXISTS idx_userId ON events(userId);
	`

	_, err = DB.Exec(createIndex)
	if err != nil {
		panic("Failed to create index for event user id")
	}
}
