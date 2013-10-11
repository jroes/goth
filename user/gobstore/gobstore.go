package gobstore

import (
	"bytes"
	"encoding/gob"
	"github.com/jroes/goth/user"
	"io/ioutil"
	"os"
	"regexp"
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

func (store UserGobStore) FindByEmail(email string) (*user.User, error) {
	emailHash := user.GenerateHash(email)
	return store.FindByHash(emailHash)
}

func (store UserGobStore) FindByHash(emailHash string) (*user.User, error) {
	// Don't trust the hash string, sanitize
	reg := regexp.MustCompile("/[^A-Za-z0-9]+/")
	safeEmailHash := reg.ReplaceAllString(emailHash, "")
	userGob, err := ioutil.ReadFile(store.Path + safeEmailHash + ".gob")
	if err != nil {
		return nil, err
	}

	userGobBuf := bytes.NewBuffer(userGob)
	decoder := gob.NewDecoder(userGobBuf)
	user := user.User{}
	err = decoder.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (store UserGobStore) Save(user user.User) error {
	emailHash := user.EmailHash()
	userGobBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(userGobBuf)
	encoder.Encode(user)
	return ioutil.WriteFile(store.Path+emailHash+".gob", userGobBuf.Bytes(), 0600)
}
