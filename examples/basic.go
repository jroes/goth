package main

import (
	"fmt"
	"net/http"
)

import (
	"github.com/jroes/goth"
)

var authHandler = goth.DefaultAuthHandler

func main() {
	http.Handle("/auth/", authHandler)
	http.HandleFunc("/", helloUserHandler)

	// Please use ListenAndServeTLS in production.
	http.ListenAndServe(":8080", nil)
}

func helloUserHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := authHandler.CurrentUser(r)
	if ok {
		fmt.Fprintf(w, "Hello, %s!", currentUser.Email)
	} else {
		fmt.Fprintf(w, "Hello, guest!")
	}
}
