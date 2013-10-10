package goth

import (
	"github.com/jroes/goth"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestSignupPostRedirects(t *testing.T) {
	ts, _ := setup()
	resp, err := http.PostForm(ts.URL+"/auth/sign_up",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		t.Errorf("Error posting to sign up route: %v", err)
	}
	if resp.StatusCode != 301 {
		t.Errorf("Expected 301 redirect after posting to sign up route. Got %d instead.", resp.StatusCode)
	}
}
