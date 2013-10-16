package goth

import (
	"fmt"
	"html/template"
	"net/http"
)

import (
	gu "github.com/jroes/goth/user"
)

// SignInHandler validates email and password parameters in an HTTP request
// against the UserStore. If the provided parameters are valid, a session will
// be created for the user and an HTTP redirect will be returned to the
// AfterSigninPath.
func (handler AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler.signInShowHandler(w, r)
	} else if r.Method == "POST" {
		handler.signInCreateHandler(w, r)
	} else {
		fmt.Println("Got a signin action with HTTP method %s, what the heck is that?", r.Method)
		http.NotFound(w, r)
	}
}

func (handler AuthHandler) signInShowHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(handler.TemplatePath + "sign_in.html")
	if err != nil {
		panic(err)
		return
	}
	t.Execute(w, handler)
}

func (handler AuthHandler) signInCreateHandler(w http.ResponseWriter, r *http.Request) {
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

// SignUpHandler handles both GET and POST requests. With a GET request, it
// renders a sign up template form that will POST to the same route. With a POST
// request, it creates a user via the UserStore with the specified email and
// password parameters. After successfully creating a User, it will redirect to
// the AfterSignupPath.
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

// SignOutHandler instructs the browser to clear the session and redirects
// the client to the AfterSignoutPath.
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
