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
	"os"
	"text/tabwriter"

	"github.com/lukas8219/gcsv/storage"
	"github.com/lukas8219/gcsv/util"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans a sheet into Stdout",
	Long:  `Scans the sheet and print the entire results into the Terminal`,
	Run:   scan,
}

func scan(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		log.Fatalln("One and only one argument is needed. Please the link or the name of the saved sheet")
	}
	columns, err := cmd.Flags().GetInt("endColumn")
	if err != nil {
		log.Fatalln(err)
	}

	end, err := cmd.Flags().GetInt("startColumn")
	if err != nil {
		log.Fatalln(err)
	}

	startChar := getColumChar(end)
	endChar := getColumChar(columns)
	batchSize, err := cmd.Flags().GetInt("size")
	if err != nil {
		log.Fatal(err)
	}

	start, err := cmd.Flags().GetInt("start")
	if err != nil {
		log.Fatalln(err)
	}

	svr := getSpreadSheetsService()
	storage := storage.GetStorage()

	selectedSheet := util.ParseLink(args[0])
	favorite, err := storage.Get(selectedSheet)
	if err == nil {
		selectedSheet = favorite
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Searching for Sheet with ID: ", selectedSheet)

	count := 0
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 8, '\t', tabwriter.AlignRight)
	defer w.Flush()

	result := make(chan interface{})

	go func(count *int) {
		for {
			columnRanges := fmt.Sprintf("%s%d:%s%d", startChar, start, endChar, batchSize)
			sheet, err := svr.Values.Get(selectedSheet, columnRanges).Do()
			if err != nil {
				log.Fatalln(err)
			}

			if len(sheet.Values) == 0 {
				break
			}

			for _, val := range sheet.Values {
				for _, entry := range val {
					result <- entry
					result <- "\t"
				}
				result <- "\n"
			}
			start = batchSize
			batchSize *= 2
		}
		close(result)
	}(&count)
	for e := range result {
		fmt.Fprint(w, e)
	}

}

func getColumChar(columns int) string {
	return string('a' + (columns - 1))
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().Int("endColumn", 5, "Defines the End column")
	scanCmd.Flags().Int("startColumn", 1, "Define the start column")

	scanCmd.Flags().Int("start", 1, "Defines the starting row")
	scanCmd.Flags().Int("size", 10, "Define the BatchSize")

}
