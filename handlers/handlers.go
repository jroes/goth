package handlers

import (
	"fmt"
	"github.com/jroes/goth/user"
	"github.com/jroes/goth/user/gobstore"
	"html/template"
	"net/http"
)

// Need to move this somewhere better
type GothConfig struct {
	PathPrefix string
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var userStore user.UserStore
	userStore = gobstore.NewUserGobStore("users/")
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := userStore.Find(email)
	if err != nil {
		fmt.Fprintf(w, "Sorry, couldn't find a user with that email address.\n")
		return
	}

	err = user.HasPassword(password)
	if err != nil {
		fmt.Fprintf(w, "Looks like you have the wrong password!\n")
		return
	}
	fmt.Fprintf(w, "Looks like you have the matching password!\n")
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		signUpShowHandler(w, r)
	} else if r.Method == "POST" {
		signUpCreateHandler(w, r)
	} else {
		fmt.Println("Got a signup action with HTTP method %s, what the heck is that?", r.Method)
		http.NotFound(w, r)
	}
}

func signUpShowHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/sign_up.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v", err)
		return
	}
	t.Execute(w, GothConfig{PathPrefix: "/auth/"})
}

func signUpCreateHandler(w http.ResponseWriter, r *http.Request) {
	var userStore user.UserStore
	userStore = gobstore.NewUserGobStore("users/")
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := user.New(email, password)
	err := userStore.Save(*user)
	if err != nil {
		fmt.Fprintf(w, "Had trouble creating %s. Error: %v\n", email, err)
		return
	}

	fmt.Fprintf(w, "Your user account has been created!\n")
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
}
