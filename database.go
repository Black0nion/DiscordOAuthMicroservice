package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var db *sql.DB

func connectToDatabase() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file!")
		return
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	databaseIP := os.Getenv("DB_IP")
	databaseName := os.Getenv("DB_NAME")

	connStr := "postgresql://%s:%s@%s/%s?sslmode=disable"
	connStr = fmt.Sprintf(connStr, username, password, databaseIP, databaseName)

	var err error
	// Connect to database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateSession(accessToken string, refreshToken string, expiresAt time.Time) (string, error) {
	var err error

	var result *sql.Rows
	result, err = db.Query("SELECT session_id FROM sessions WHERE access_token = $1 AND refresh_token = $2", accessToken, refreshToken)
	if err != nil {
		return "", err
	}

	// If the result is not empty, then we don't need to insert a new row
	if result.Next() {
		var id string
		err = result.Scan(&id)
		if err != nil {
			return "", err
		}

		return id, nil
	}

	id := GenerateSessionID(db)

	_, err = db.Exec("INSERT INTO sessions (session_id, access_token, refresh_token, expires_at) VALUES ($1, $2, $3, $4)", id, accessToken, refreshToken, expiresAt.UnixMilli())

	return id, err
}
