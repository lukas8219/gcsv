package cmd

import (
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const TOKEN_FILE_PATH = "token.json"
const SECRETS_FILE_PATH = "secret.json"

func authenticate() (*oauth2.Config, *oauth2.Token, error) {
	file, err := ioutil.ReadFile(SECRETS_FILE_PATH)

	if err != nil {
		return nil, nil, err
	}

	configs, err := google.ConfigFromJSON(file, SCOPE)
	if err != nil {
		return nil, nil, err
	}

	token, err := tokenFromFile()
	return configs, token, err
}
