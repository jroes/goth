package goth

import (
	"github.com/jroes/goth"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup() (*httptest.Server, goth.AuthHandler) {
	authHandler := goth.AuthHandler{RoutePath: "/auth/", TemplatePath: "tmpl/", AfterSignupURL: "/", AfterSigninURL: "/"}
	return httptest.NewServer(authHandler), authHandler
}

func TestSignUpGetRendersWithPathPrefix(t *testing.T) {
	ts, authHandler := setup()
	resp, err := http.Get(ts.URL + "/auth/sign_up")
	if err != nil {
		t.Errorf("Error visiting sign up page: %v", err)
	}
	if resp.StatusCode > 400 {
		t.Errorf("Error status when visiting sign up page: %d", resp.StatusCode)
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(contents), authHandler.RoutePath) {
		t.Errorf("Did not find configured RoutePath (%s) within the response body.", authHandler.RoutePath)
	}
}

func TestSignupPostHandled(t *testing.T) {
	ts, _ := setup()
	resp, err := http.Get(ts.URL + "/auth/sign_up")
	if err != nil {
		t.Errorf("Error visiting sign up page: %v", err)
	}
	if resp.StatusCode > 400 {
		t.Errorf("Error status when visiting sign up page: %d", resp.StatusCode)
	}
}
