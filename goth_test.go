package goth_test

import (
	"fmt"
	"github.com/jroes/goth/user/gobstore"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

import "github.com/jroes/goth"

var testMux = http.NewServeMux()
var testUserStorePath = "/tmp/users/"

func makeHelloUserHandler(authHandler goth.AuthHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := authHandler.CurrentUser(r)
		fmt.Fprintf(w, "Hello, %s!", currentUser.Email)
	}
}

func cleanup() {
	os.RemoveAll(testUserStorePath)
}

func setup() (*httptest.Server, goth.AuthHandler) {
	// Reset test muxer each run
	testMux = http.NewServeMux()
	authHandler := goth.DefaultAuthHandler
	authHandler.UserStore = gobstore.NewUserGobStore(testUserStorePath)
	testMux.HandleFunc(authHandler.RoutePath, authHandler.ServeHTTP)
	testMux.HandleFunc("/", makeHelloUserHandler(authHandler))
	return httptest.NewServer(testMux), authHandler
}

func TestSignupGetRendersWithPathPrefix(t *testing.T) {
	ts, authHandler := setup()
	defer ts.Close()
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

func TestSignupPostLogsInAndRedirects(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()
	defer cleanup()
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	resp, err := client.PostForm(ts.URL+"/auth/sign_up",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		t.Errorf("Error posting to sign up route: %v", err)
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "jon@example.com") {
		t.Errorf("Expected response to contain jon@example.com: %s", contents)
	}
}

func TestSigninPostLogsInAndRedirects(t *testing.T) {
	ts, _ := setup()
	defer ts.Close()
	defer cleanup()
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	resp, err := client.PostForm(ts.URL+"/auth/sign_in",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		t.Errorf("Error posting to sign up route: %v", err)
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "jon@example.com") {
		t.Errorf("Expected response to contain jon@example.com: %s", contents)
	}
}
