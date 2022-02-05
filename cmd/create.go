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
	"log"

	"github.com/lukas8219/gcsv/storage"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			log.Fatal("One and only one argument is required [Title]")
		}

		svc := getSpreadSheetsService()

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}

		validateName(name)

		res, err := svc.Create(&sheets.Spreadsheet{
			Properties: &sheets.SpreadsheetProperties{
				Title: args[0],
			},
		}).Do()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Sheet '%s' saved successfully\n", args[0])
		log.Printf("Saving it into Favorite")

		database := storage.GetStorage()

		database.Store(storage.FavoriteSheet{
			Name: name,
			ID:   res.SpreadsheetId,
		})

	},
}

func validateName(name string) {
	storage := storage.GetStorage()
	_, err := storage.Get(name)
	if err == nil {
		log.Fatal("Name already exists. Please give a different name to save it as favorite or remove the other")
	}
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("name", "", "name")
	createCmd.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
