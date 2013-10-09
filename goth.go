package goth

import (
	"fmt"
	"github.com/jroes/goth/handlers"
	"net/http"
	"regexp"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	actionRegexp := regexp.MustCompile(".*\\/(.*)")
	actionMatches := actionRegexp.FindStringSubmatch(r.URL.Path)
	if actionMatches == nil || len(actionMatches) != 2 {
		fmt.Printf("actionMatches was %q for %s", actionMatches, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	action := actionMatches[1]

	if action == "sign_in" {
		handlers.SignInHandler(w, r)
	} else if action == "sign_out" {
		handlers.SignOutHandler(w, r)
	} else if action == "sign_up" {
		handlers.SignUpHandler(w, r)
	}

	http.NotFound(w, r)
}
