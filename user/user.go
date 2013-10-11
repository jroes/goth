package user

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha1"
	"encoding/base64"
)

type User struct {
	Email        string
	PasswordHash []byte
}

// HasPassword uses bcrypt to compare the supplied plaintext password with
// the stored password hash for the user.
func (user *User) HasPassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}

// EmailHash returns hashed version of the user's email address.
func (user *User) EmailHash() string {
	return GenerateHash(user.Email)
}

func GenerateHash(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// New creates a new User type with an email address and plaintext password.
// The password is immediately hashed with bcrypt and stored in User.PasswordHash.
func New(email string, password string) *User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return &User{email, hashedPassword}
}

// The UserStore interface is implemented to persist and retrieve User types.
type UserStore interface {
	FindByEmail(string) (*User, error)
	FindByHash(string) (*User, error)
	Save(User) error
}
