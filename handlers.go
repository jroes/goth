package goth

import (
	"fmt"
	gu "github.com/jroes/goth/user"
	"html/template"
	"net/http"
)

func (handler AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := handler.UserStore.FindByEmail(email)
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

func (handler AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler.signUpShowHandler(w, r)
	} else if r.Method == "POST" {
		handler.signUpCreateHandler(w, r)
	} else {
		fmt.Println("Got a signup action with HTTP method %s, what the heck is that?", r.Method)
		http.NotFound(w, r)
	}
}

func (handler AuthHandler) signUpShowHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(handler.TemplatePath + "sign_up.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v", err)
		return
	}
	t.Execute(w, handler)
}

func (handler AuthHandler) signUpCreateHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user := gu.New(email, password)
	err := handler.UserStore.Save(*user)
	if err != nil {
		fmt.Fprintf(w, "Had trouble creating %s. Error: %v\n", email, err)
		return
	}

	session, err := handler.SessionStore.Get(r, "goth-session")
	if err != nil {
		panic(err)
	}
	session.Values["identifier"] = user.EmailHash()
	session.Save(r, w)

	http.Redirect(w, r, handler.AfterSignupPath, http.StatusFound)
}

func (handler AuthHandler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
}
