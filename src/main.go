package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
		return
	}

	connectToDatabase() // <--- starts the database connection
	handleAuth()        // <--- starts everything we need with the discord API
}
