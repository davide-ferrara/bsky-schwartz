package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNoRow = fmt.Errorf("no rows in result set")

var db *sql.DB

func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("could not create tables: %w", err)
	}

	return nil
}

func createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		handle TEXT PRIMARY KEY,
		did TEXT NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	weightsTable := `
	CREATE TABLE IF NOT EXISTS weights (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_handle TEXT NOT NULL,
		value_id TEXT NOT NULL,
		weight REAL NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_handle) REFERENCES users(handle) ON DELETE CASCADE,
		UNIQUE(user_handle, value_id)
	);
	`

	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	if _, err := db.Exec(weightsTable); err != nil {
		return err
	}

	return nil
}

func SaveUser(handle, did string) error {
	_, err := db.Exec(`
		INSERT INTO users (handle, did, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(handle) DO UPDATE SET did = excluded.did, updated_at = excluded.updated_at
	`, handle, did, time.Now())
	return err
}

func SaveWeight(userHandle, valueID string, weight float64) error {
	_, err := db.Exec(`
		INSERT INTO weights (user_handle, value_id, weight, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_handle, value_id) DO UPDATE SET weight = excluded.weight, updated_at = excluded.updated_at
	`, userHandle, valueID, weight, time.Now())
	return err
}

func SaveWeights(userHandle string, weights map[string]float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO weights (user_handle, value_id, weight, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_handle, value_id) DO UPDATE SET weight = excluded.weight, updated_at = excluded.updated_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for valueID, weight := range weights {
		if _, err := stmt.Exec(userHandle, valueID, weight, time.Now()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetUserFromDB(handle string) (string, map[string]float64, error) {
	var did string
	err := db.QueryRow("SELECT did FROM users WHERE handle = ?", handle).Scan(&did)
	if err != nil {
		return "", nil, fmt.Errorf("user not found: %w", err)
	}

	rows, err := db.Query("SELECT value_id, weight FROM weights WHERE user_handle = ?", handle)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	weights := make(map[string]float64)
	for rows.Next() {
		var valueID string
		var weight float64
		if err := rows.Scan(&valueID, &weight); err != nil {
			return "", nil, err
		}
		weights[valueID] = weight
	}

	return did, weights, nil
}

func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
