package goth

import (
	"fmt"
	"net/http"
	"regexp"
)

type AuthHandler struct {
	RoutePath      string
	TemplatePath   string
	AfterSignupURL string
	AfterSigninURL string
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
