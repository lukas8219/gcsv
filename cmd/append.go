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
	"fmt"
	"log"
	"strings"

	"github.com/lukas8219/gcsv/storage"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := getSpreadSheetsService()

		if len(args) != 1 {
			log.Fatal("One and only one argument is required!")
		}

		delimiter, err := cmd.Flags().GetString("d")
		if err != nil {
			log.Fatal(err)
		}

		storage := storage.GetStorage()

		input := args[0]
		entry := strings.Split(input, delimiter)

		values := make([][]interface{}, 1)
		values[0] = make([]interface{}, len(entry))
		v := values[0]
		for i, val := range entry {
			v[i] = val
		}

		value := &sheets.ValueRange{
			Values: values,
		}

		param, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}

		id, err := storage.Get(param)
		if err != nil {
			log.Fatal(err)
		}

		startChar := 'A'
		endChar := getEndChar(entry)

		rangeId := fmt.Sprintf("%c1:%s1", startChar, endChar)

		_, err = svc.Values.Append(id, rangeId, value).ValueInputOption("RAW").Do()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Success!")
	},
}

func getEndChar(entry []string) string {
	return string('a' + len(entry))
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().String("d", ",", "Delimiter")
	appendCmd.Flags().String("name", "", "Name")

	appendCmd.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
