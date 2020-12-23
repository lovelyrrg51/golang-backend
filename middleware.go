package main

import (
	"fmt"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check the auth token exists in the request
		keys := []string{"auth_token", "user_name"}

		userName := ""
		authToken := ""

		// check URL query parameters eg. https://domain.com/home?auth_token=xxxxxxxx
		getKeys := retreiveGetParameters(keys, r)
		if val, ok := getKeys["user_name"]; ok {
			userName = val
		}
		log.Println(userName)
		if val, ok := getKeys["auth_token"]; ok {
			authToken = val
		}
		log.Println(authToken)

		// check POST form parameters from application/x-www-form-urlencoded content type
		postKeys := retreivePostParameters(keys, r)
		if val, ok := postKeys["user_name"]; ok {
			userName = val
		}
		if val, ok := postKeys["auth_token"]; ok {
			authToken = val
		}

		if userName != "" && authToken != "" {

			db := Dbconn()

			var user User
			// check that the database has that auth token for that username
			err := db.QueryRow("SELECT name, id, authToken FROM users WHERE user_name = ?", userName).Scan(&user.Name, &user.ID, &user.AuthToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error: Authentication request returned an error: %s", err), http.StatusInternalServerError)
				return
			}

			// check the db had the correct authtoken for user with that name
			if user.AuthToken == authToken {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Authentication failed, auth token was not found in the database", http.StatusForbidden)
			return
		}

		http.Error(w, "Authentication Token was not found in request", http.StatusForbidden)
	})
}

func retreiveGetParameters(keys []string, r *http.Request) map[string]string {
	m := make(map[string]string)
	//iterate over the keys
	for _, k := range keys {
		log.Println(k)
		v := r.URL.Query().Get(k)
		if v == "" {
			continue
		}
		m[k] = v
	}
	return m
}

func retreivePostParameters(keys []string, r *http.Request) map[string]string {
	m := make(map[string]string)

	//iterate over the keys
	for k := range m {
		v := r.FormValue(k)
		if v == "" {
			continue
		}
		m[k] = v
	}

	return m
}

type User struct {
	Name      string
	AuthToken string
	ID        int64
}
