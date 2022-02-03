package cmd

import (
	"context"
	"log"
	"net/http"
)

func getClient() *http.Client {
	config, token, err := authenticate()
	if err != nil {
		log.Fatalln("You are not authenticated.")
	}

	return config.Client(context.Background(), token)
}
