package cmd

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getSpreadSheetsService() *sheets.SpreadsheetsService {
	client := getClient()
	ctx := context.Background()
	svr, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)

	}
	return svr.Spreadsheets
}

func getClient() *http.Client {
	config, token, err := authenticate()
	if err != nil {
		log.Fatalln("You are not authenticated.")
	}

	return config.Client(context.Background(), token)
}
