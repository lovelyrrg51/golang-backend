package main

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Do something with the request

	//then return
	w.Write([]byte("Success"))
}
