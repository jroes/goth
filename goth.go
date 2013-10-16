// Package goth provides an authentication system for Go web apps.
package goth

import (
	"fmt"
	"net/http"
	"regexp"
)

import (
	"github.com/gorilla/sessions"
	. "github.com/jroes/goth/user"
	"github.com/jroes/goth/user/gobstore"
)

type AuthHandler struct {
	// Where to mount URLs for authentication (e.g. signup, signin)
	RoutePath string
	// Where on disk HTML templates are stored for authentication pages.
	TemplatePath string
	// Where to redirect the user after various authentication operations.
	AfterSignupPath  string
	AfterSigninPath  string
	AfterSignoutPath string
	// This should be set to a string of characters used to encrypt and sign
	// sessions. It should be kept private from any source repositories. You could
	// use os.Getenv() for this and store it in an environment variable.
	SessionSecret string
	SessionStore  *sessions.CookieStore
	UserStore     UserStore
}

var DefaultAuthHandler = AuthHandler{
	RoutePath:        "/auth/",
	TemplatePath:     "tmpl/",
	AfterSignupPath:  "/",
	AfterSigninPath:  "/",
	AfterSignoutPath: "/",
	SessionSecret:    "change-me-please",
	SessionStore:     sessions.NewCookieStore([]byte("change-me-please")),
	UserStore:        gobstore.NewUserGobStore("users/"),
}

// ServeHTTP implements the http.Handler interface to delegate
// authentication-related routing to the proper handler.
func (handler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actionRegexp := regexp.MustCompile(".*\\/(.*)")
	actionMatches := actionRegexp.FindStringSubmatch(r.URL.Path)
	if actionMatches == nil || len(actionMatches) != 2 {
		fmt.Printf("Unexpected request: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	action := actionMatches[1]

	if action == "sign_in" {
		handler.SignInHandler(w, r)
		return
	} else if action == "sign_out" {
		handler.SignOutHandler(w, r)
		return
	} else if action == "sign_up" {
		handler.SignUpHandler(w, r)
		return
	}

	http.NotFound(w, r)
}

// CurrentUser retrieves the User object for the currently logged in user based
// on the request. The first return value is the User object, the second is
// true if a user is logged in. If a user is not logged in, the first return
// value will be an empty User, and the second return value will be false.
func (handler AuthHandler) CurrentUser(r *http.Request) (*User, bool) {
	session, err := handler.SessionStore.Get(r, "goth-session")
	if err != nil {
		panic(err)
	}
	emailHash, ok := session.Values["identifier"]
	if ok {
		user, err := handler.UserStore.FindByHash(emailHash.(string))
		if err != nil {
			panic(fmt.Errorf("Couldn't find user with identifier %s in user store.", emailHash.(string)))
		}
		return user, true
	}
	return &User{}, false
}
