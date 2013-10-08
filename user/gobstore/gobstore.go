package gobstore

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"os"
)

// The UserGobStore implements the UserStore interface to store Users on disk
// encoded in the gob format.
type UserGobStore struct {
	Path string
}

func NewUserGobStore(path string) *UserGobStore {
	store := UserGobStore{path}
	err := os.MkdirAll(path, 0700)
	if err != nil {
		panic(err)
	}
	return &store
}

func (store UserGobStore) Find(email string) (*User, error) {
	emailSha := generateHash(email)
	userGob, err := ioutil.ReadFile(store.Path + emailSha + ".gob")
	if err != nil {
		return nil, err
	}

	userGobBuf := bytes.NewBuffer(userGob)
	decoder := gob.NewDecoder(userGobBuf)
	user := User{}
	err = decoder.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (store UserGobStore) Save(user User) error {
	emailSha := generateHash(user.Email)
	userGobBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(userGobBuf)
	encoder.Encode(user)
	return ioutil.WriteFile(store.Path+emailSha+".gob", userGobBuf.Bytes(), 0600)
}

func generateHash(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
