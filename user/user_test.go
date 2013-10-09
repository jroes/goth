package user

import (
	"github.com/jroes/goth/user"
	"testing"
)

func TestHasPassword(t *testing.T) {
	user := user.New("test@example.com", "password")
	if err := user.HasPassword("password"); err != nil {
		t.Errorf("Failed to recognize valid password for user, %v.\n", err)
	}
}

func TestSetsPasswordHash(t *testing.T) {
	user := user.New("test@example.com", "password")
	if user.PasswordHash == nil {
		t.Errorf("No PasswordHash was created.")
	}
}
