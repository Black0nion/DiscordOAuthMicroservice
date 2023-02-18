package main

import (
	"context"
	"encoding/json"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

// This is the state key used for security, sent in login, validated in callback.
// For this example we keep it simple and hardcode a string
// but in real apps you must provide a proper function that generates a state.
var state = "random"

var corsHeaderValue string

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", corsHeaderValue)
}

func handleAuth() {
	corsHeaderValue = GetEnvOrDefault("CORS_ORIGIN", "*")

	conf := &oauth2.Config{
		RedirectURL: GetEnv("REDIRECT_URL"),
		// This next 2 lines must be edited before running this.
		ClientID:     GetEnv("CLIENT_ID"),
		ClientSecret: GetEnv("CLIENT_SECRET"),
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
	http.HandleFunc("/auth/exchange_code", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.FormValue("state") != state {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("State does not match."))
			return
		}

		// Step 3: We exchange the code we got for an access token
		// Then we can use the access token to do actions, limited to scopes we requested
		token, err := conf.Exchange(context.Background(), r.FormValue("code"))

		if err != nil {
			// check if the error is of type *oauth2.RetrieveError
			if oauthError, ok := err.(*oauth2.RetrieveError); ok {
				w.WriteHeader(oauthError.Response.StatusCode)

				// load body from error and parse it as JSON
				var body map[string]any
				err := json.Unmarshal(oauthError.Body, &body)
				if err == nil {
					i := body["error_description"]

					bytes := []byte(i.(string))
					_, _ = w.Write(bytes)

					return
				}
			}

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
		_, _ = w.Write([]byte(sid))
	})

	handleUser()
	handleGuilds()
	port := ":" + GetEnv("PORT")
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
