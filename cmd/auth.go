/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/lukas8219/gcsv/storage"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const SCOPE = "https://www.googleapis.com/auth/spreadsheets"

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Log-in with your Google Account",
	Long:  `Presents a Google Authentication mechanism for log-in with Google Account and get access to Sheets`,
	Run:   auth,
}

func auth(cmd *cobra.Command, args []string) {

	configs, token, err := authenticate()
	if err != nil {
		log.Println("Error occurred when trying to fetch Token locally. Retrieving another one\t\t", err)
		token = getTokenFromWeb(configs)
		saveToken(token)
	} else {
		log.Println("You're already authenticated!")
	}
	_ = configs.Client(context.Background(), token)
}

func saveToken(token *oauth2.Token) {
	log.Printf("Saving credentials in %s", storage.GetTokenFilePath())
	f, err := storage.CreateOrWriteTokenFile()
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
	log.Println("Successfully saved the Token into")
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Println("Insert the Authentication Code Here: ")

	//Need to properly find a way to execute the default browser
	err := exec.Command("google-chrome", authURL).Run()

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	log.Println("Sucessfully authenticated!")
	return tok
}

func tokenFromFile() (*oauth2.Token, error) {
	f, err := storage.ReadTokenFile()
	defer f.Close()

	if err != nil {
		return nil, err
	}

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
