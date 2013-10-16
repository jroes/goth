package goth_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

import (
	"github.com/jroes/goth"
	"github.com/jroes/goth/user"
	"github.com/jroes/goth/user/gobstore"
)

var testMux = http.NewServeMux()
var testUserStorePath = "/tmp/users/"

func cleanup(ts *httptest.Server) {
	os.RemoveAll(testUserStorePath)
	ts.Close()
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

func makeHelloUserHandler(authHandler goth.AuthHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := authHandler.CurrentUser(r)
		if ok {
			fmt.Fprintf(w, "Hello, %s!", currentUser.Email)
		} else {
			fmt.Fprintf(w, "Hello, guest!")
		}
	}
}

func createAndSignInUser(ts *httptest.Server, authHandler goth.AuthHandler, client http.Client) {
	user := user.New("jon@example.com", "password")
	authHandler.UserStore.Save(*user)
	resp, err := client.PostForm(ts.URL+"/auth/sign_in",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		panic(err)
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "jon@example.com") {
		panic(fmt.Sprintf("Expected response to contain jon@example.com: %s", contents))
	}
}

func TestSignupGetRendersWithRoutePath(t *testing.T) {
	ts, authHandler := setup()
	defer cleanup(ts)
	resp, err := http.Get(ts.URL + "/auth/sign_up")
	if err != nil {
		t.Errorf("Error visiting sign up page: %v", err)
		return
	}
	if resp.StatusCode > 400 {
		t.Errorf("Error status when visiting sign up page: %d", resp.StatusCode)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(contents), authHandler.RoutePath) {
		t.Errorf("Did not find configured RoutePath (%s) within the response body.", authHandler.RoutePath)
		return
	}
}

func TestSignupPostLogsInAndRedirects(t *testing.T) {
	ts, _ := setup()
	defer cleanup(ts)
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	resp, err := client.PostForm(ts.URL+"/auth/sign_up",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		t.Errorf("Error posting to sign up route: %v", err)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "jon@example.com") {
		t.Errorf("Expected response to contain jon@example.com: %s", contents)
		return
	}
}

func TestSigninGetRendersWithRoutePath(t *testing.T) {
	ts, authHandler := setup()
	defer cleanup(ts)
	resp, err := http.Get(ts.URL + "/auth/sign_in")
	if err != nil {
		t.Errorf("Error visiting sign in page: %v", err)
		return
	}
	if resp.StatusCode > 400 {
		t.Errorf("Error status when visiting sign in page: %d", resp.StatusCode)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(contents), authHandler.RoutePath) {
		t.Errorf("Did not find configured RoutePath (%s) within the response body.", authHandler.RoutePath)
		return
	}
}

func TestSigninPostLogsInAndRedirects(t *testing.T) {
	ts, authHandler := setup()
	defer cleanup(ts)
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	user := user.New("jon@example.com", "password")
	authHandler.UserStore.Save(*user)
	resp, err := client.PostForm(ts.URL+"/auth/sign_in",
		url.Values{"email": {"jon@example.com"}, "password": {"password"}})
	if err != nil {
		t.Errorf("Error posting to sign in route: %v", err)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "jon@example.com") {
		t.Errorf("Expected response to contain jon@example.com: %s", contents)
		return
	}
}

func TestSignoutPostRedirectsAndDoesntKnowYou(t *testing.T) {
	ts, authHandler := setup()
	defer cleanup(ts)
	client := &http.Client{}
	client.Jar, _ = cookiejar.New(nil)
	createAndSignInUser(ts, authHandler, *client)
	resp, err := client.PostForm(ts.URL+"/auth/sign_out", url.Values{})
	if err != nil {
		t.Errorf("Error posting to sign out route: %v", err)
		return
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if !strings.Contains(string(contents), "guest") {
		t.Errorf("Expected response to contain guest: %s", contents)
		return
	}
}
