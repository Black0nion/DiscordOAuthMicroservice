package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const targetStringLength = 69
const leftLimit1 = 'A'
const leftLimit2 = 'a'
const leftLimit3 = '0'
const rightLimit1 = 'Z'
const rightLimit2 = 'z'
const rightLimit3 = '9'

func GenerateSessionID(db *sql.DB) string {
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	var generatedIDBuilder strings.Builder
	for i := 0; i < targetStringLength; i++ {
		r := randSource.Intn(rightLimit2+1) + leftLimit3
		if r >= leftLimit1 && r <= rightLimit1 || r >= leftLimit2 && r <= rightLimit2 || r >= leftLimit3 && r <= rightLimit3 {
			generatedIDBuilder.WriteRune(rune(r))
		} else {
			i--
		}
	}

	if rows, _ := db.Query("SELECT session_id FROM sessions WHERE session_id = $1", generatedIDBuilder.String()); rows.Next() {
		return GenerateSessionID(db)
	}

	return generatedIDBuilder.String()
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is empty!", key)
	}
	return value
}

func GetEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
