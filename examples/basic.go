package main

import (
	"fmt"
	"html/template"
	"net/http"
)

import (
	"github.com/jroes/goth"
)

var authHandler = goth.DefaultAuthHandler

type WelcomeData struct {
	SignupURL string
	SigninURL string
	Greeting  string
}

func main() {
	http.Handle("/auth/", authHandler)
	http.HandleFunc("/", helloUserHandler)

	// Please use ListenAndServeTLS in production.
	http.ListenAndServe(":8080", nil)
}

func helloUserHandler(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := authHandler.CurrentUser(r)
	var greeting string
	if ok {
		greeting = fmt.Sprintf("Hello, %s!", currentUser.Email)
	} else {
		greeting = "Hello, guest!"
	}

	welcome := WelcomeData{
		SignupURL: authHandler.RoutePath + "sign_up",
		SigninURL: authHandler.RoutePath + "sign_in",
		Greeting:  greeting,
	}
	t := template.Must(template.New("welcome").Parse(`
<html>
{{.Greeting}}
<ul>
  <li><a href='{{.SigninURL}}'>Sign in</a></li>
  <li><a href='{{.SignupURL}}'>Sign up</a></li>
</ul>
</html>`))
	t.Execute(w, welcome)
}
