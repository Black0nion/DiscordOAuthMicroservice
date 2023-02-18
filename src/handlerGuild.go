package main

import "net/http"

func handleGuilds() {
	http.HandleFunc("/users/@me/guilds", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == "GET" {
			//get the json data from the request
			w.WriteHeader(http.StatusOK)
			var body []byte
			_, _ = w.Write(body)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
