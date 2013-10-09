package goth

import (
	"github.com/jroes/goth"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(goth.AuthHandler))
}

func TestSignUpShowHandled(t *testing.T) {
	ts := setup()
	resp, err := http.Get(ts.URL + "/auth/sign_up")
	if err != nil {
		t.Errorf("Error visiting sign up page: %v", err)
	}
	if resp.StatusCode > 400 {
		t.Errorf("Error status when visiting sign up page: %d", resp.StatusCode)
	}
}
