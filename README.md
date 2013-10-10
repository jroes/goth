# Goth

Goth is a web authentication system written in Go. With Goth, you get out-of-
the-box user sign in, sign up, and sign out functionality to kick off building
your web app.

## Installation
```
go get github.com/jroes/goth
```

## Usage
The following example registers a handler for an admin section that would be
able to retrieve information about the currently logged in user, and a pattern
that will handle the standard sign in, sign out, and sign up functionality.

```go
package main

import (
    "fmt"
    "github.com/jroes/goth"
    "net/http"
)

func main() {
    authHandler := goth.AuthHandler{RoutePath: "/auth/", TemplatePath: "tmpl/", AfterSignupURL: "/", AfterSigninURL: "/"}
    http.HandleFunc("/admin/", adminHandler)
    http.Handle("/auth/", authHandler)

    // Please use ListenAndServeTLS in production.
    http.ListenAndServe(":8080", nil)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
    if user, err := goth.CurrentUser(r); err != nil {
        fmt.Fprint(w, "User not logged in, please authenticate before visiting this page.")
        return
    }

    fmt.Fprint(w, "Hello, %s", user.email)
}
```
