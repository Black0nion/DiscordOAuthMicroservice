package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"io/ioutil"
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
		RedirectURL: "http://localhost:3000/auth/callback",
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
			w.Write([]byte(err.Error()))
			return
		}

		// Step 4: Use the access token, here we use it to get the logged in user's info.
		res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write([]byte(res.Status))
			}
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(body)
	})

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
