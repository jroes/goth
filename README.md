# Goth
[![Build Status](https://drone.io/github.com/jroes/goth/status.png)](https://drone.io/github.com/jroes/goth/latest)

Goth is a web authentication system written in Go. With Goth, you get out-of-
the-box user sign in, sign up, and sign out functionality to kick off building
your web app.

## Disclaimer
Not ready for production yet! Use at your own risk. File bugs, pitch in, etc.

## Installation
```
go get github.com/jroes/goth
```

## Usage
The following example mounts goth's authentication handler at `/auth/`, and
registers a simple handler for the root route that greets the user with the
user's e-mail address if logged in, or as a guest if not.

```go
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
```

So what does this get you exactly? You get a bunch of routes underneath `/auth/`:
 * GET http://localhost:8080/auth/sign_up - Standard sign up form
 * POST http://localhost:8080/auth/sign_up - Creates User and persists using the defined UserStore (default: UserGobStore, a store persisted in Go's gob format on disk)
 * GET http://localhost:8080/auth/sign_in - Standard sign in form
 * POST http://localhost:8080/auth/sign_in - Creates a session (default: cookie-based through [gorilla/sessions](https://github.com/gorilla/sessions))
 * POST http://localhost:8080/auth/sign_out - Expires the session (default: expires the cookie from gorilla/sessions)
