package main

func main() {
	server()     // <--- This is the line that starts the database connection
	handleAuth() // <--- This is the line that starts everything we need with the discord API
}
