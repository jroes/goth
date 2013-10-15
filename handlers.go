package goth

import (
	"fmt"
	"html/template"
	"net/http"
)

import (
	gu "github.com/jroes/goth/user"
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
		fmt.Fprintf(w, "Sorry, the password you provided does not match.\n")
		return
	}

	handler.createUserSession(r, w, user)
	http.Redirect(w, r, handler.AfterSigninPath, http.StatusFound)
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

	handler.createUserSession(r, w, user)
	http.Redirect(w, r, handler.AfterSignupPath, http.StatusFound)
}

func (handler AuthHandler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := handler.SessionStore.Get(r, "goth-session")
	if err != nil {
		panic(err)
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, handler.AfterSignoutPath, http.StatusFound)
}

func (handler AuthHandler) createUserSession(r *http.Request, w http.ResponseWriter, user *gu.User) {
	session, err := handler.SessionStore.Get(r, "goth-session")
	if err != nil {
		panic(err)
	}
	session.Values["identifier"] = user.EmailHash()
	session.Save(r, w)
}
