package utils

import (
	"database/sql"
	"os"
	"path/filepath"

	gonanoid "github.com/matoous/go-nanoid/v2"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	dbDir := "./db"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dbDir, "database.db")
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT
	);

	CREATE TABLE IF NOT EXISTS config (
		id TEXT PRIMARY KEY,
		path TEXT,
		username TEXT,
		password TEXT,
		root TEXT
	);
	`

	_, err = DB.Exec(createTablesSQL)
	return err
}

func GenerateID() string {
	id, _ := gonanoid.New(21)
	return id
}
