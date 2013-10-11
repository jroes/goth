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
	RoutePath       string
	TemplatePath    string
	AfterSignupPath string
	AfterSigninPath string
	SessionSecret   string
	SessionStore    *sessions.CookieStore
	UserStore       UserStore
}

var DefaultAuthHandler = AuthHandler{
	RoutePath:       "/auth/",
	TemplatePath:    "tmpl/",
	AfterSignupPath: "/",
	AfterSigninPath: "/",
	SessionSecret:   "change-me-please",
	SessionStore:    sessions.NewCookieStore([]byte("change-me-please")),
	UserStore:       gobstore.NewUserGobStore("users/"),
}

func (handler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	actionRegexp := regexp.MustCompile(".*\\/(.*)")
	actionMatches := actionRegexp.FindStringSubmatch(r.URL.Path)
	if actionMatches == nil || len(actionMatches) != 2 {
		fmt.Printf("actionMatches was %q for %s", actionMatches, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	action := actionMatches[1]

	if action == "sign_in" {
		handler.SignInHandler(w, r)
	} else if action == "sign_out" {
		handler.SignOutHandler(w, r)
	} else if action == "sign_up" {
		handler.SignUpHandler(w, r)
	}

	http.NotFound(w, r)
}

func (handler AuthHandler) CurrentUser(r *http.Request) *User {
	session, err := handler.SessionStore.Get(r, "goth-session")
	if err != nil {
		panic(err)
	}
	emailHash, ok := session.Values["identifier"]
	if ok {
		user, err := handler.UserStore.FindByHash(emailHash.(string))
		if err != nil {
			panic(fmt.Sprintf("Couldn't find user with identifier %s in user store.", emailHash.(string)))
		}
		return user
	}
	return &User{}
}
