package cmd

import (
	"bytes"
	"log"

	"github.com/lukas8219/gcsv/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const TOKEN_FILE_PATH = "token.json"

func authenticate() (*oauth2.Config, *oauth2.Token, error) {
	file, err := storage.GetSecretFile()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Readed the Secret File", file.Name())
	fbyte := bytes.NewBuffer(nil)
	fbyte.ReadFrom(file)
	configs, err := google.ConfigFromJSON(fbyte.Bytes(), SCOPE)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Reading the Token Locally")
	token, err := tokenFromFile()
	return configs, token, err
}
