package config

import (
	"database/sql"
	"fmt"

	"github.com/sergiocltn/apartment-scrapper/internal/provider"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitSQLite() error {
	var err error
	DB, err = sql.Open("sqlite", GetConfig().DBPath)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		DB.Close()
		return fmt.Errorf("failed to ping SQLite database: %v", err)
	}

	provider.InfoLogger.Println("SQLite database connection initialized")
	return nil
}

func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close SQLite database: %v", err)
		}
		provider.InfoLogger.Println("SQLite database connection closed")
	}
	return nil
}
