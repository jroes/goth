# Goth

Goth is a web authentication system written in Go. With Goth, you get out-of-
the-box user sign in, sign up, and sign out functionality to kick off building y
our web app.

## Installation
```
go get github.com/jroes/goth
```

## Usage
The following example registers a pattern that requires authentication
(`/admin/`), and a pattern that will handle the standard sign in, sign out, and
sign up functionality.

```go
package main

import (
    "github.com/jroes/goth"
    "http"
)

func main() {
    http.HandleFunc("/admin/", goth.AuthRequiredHandler(adminHandler))
    http.HandleFunc("/auth/", goth.AuthHandler)

    // You should use ListenAndServeTLS in production.
    http.ListenAndServe(":8080", nil)
}
```
