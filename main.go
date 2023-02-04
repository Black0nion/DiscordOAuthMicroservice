package main

func main() {
	connectToDatabase() // <--- starts the database connection
	handleAuth()        // <--- starts everything we need with the discord API
}
