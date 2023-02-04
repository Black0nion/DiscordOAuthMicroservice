package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

// This is the state key used for security, sent in login, validated in callback.
// For this example we keep it simple and hardcode a string
// but in real apps you must provide a proper function that generates a state.
var state = "random"

func handleAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
		return
	}

	conf := &oauth2.Config{
		RedirectURL: os.Getenv("PUBLIC_URL") + "/auth/callback",
		// This next 2 lines must be edited before running this.
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify, discord.ScopeGuilds},
		Endpoint:     discord.Endpoint,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Redirect to the OAuth 2.0 Authorization page.
		// This route could be named /login etc
		http.Redirect(w, r, conf.AuthCodeURL(state)+"&prompt=none", http.StatusTemporaryRedirect)
	})

	// Step 2: After user authenticates their accounts this callback is fired.
	// the state we sent in login is also sent back to us here
	// we have to verify it as necessary before continuing.
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != state {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("State does not match."))
			return
		}

		// Step 3: We exchange the code we got for an access token
		// Then we can use the access token to do actions, limited to scopes we requested
		token, err := conf.Exchange(context.Background(), r.FormValue("code"))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		var sid string
		sid, err = CreateSession(token.AccessToken, token.RefreshToken, token.Expiry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Step 4: We can now return the created session id to the user
		_, err = w.Write([]byte(sid))
		if err != nil {
			return
		}
	})

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
