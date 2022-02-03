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
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/lukas8219/gcsv/util"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: scan,
}

func scan(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		log.Fatalln("One and only one argument is needed. Please the link or the name of the saved sheet")
	}

	selectedSheet := util.Parse(args[0])
	log.Println("Searching for Sheet with ID: ", selectedSheet)
	//TODO parse this args. Check into cache. If not, and contains HTTP or HTTPS -> try to parse only the ID
	//Else, go with the whole arg

	client := getClient()
	ctx := context.Background()

	svr, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)
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

	//Try implementing with GoRoutines
	//Do an BinarySearch-like algorithm to find all Cells
	//Request a Batch from a fixed size column
	//Cache it
	//Duplicate the batch size, starting from the previous position
	//Repeat until no more batches return

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, '\t', 0)
	defer w.Flush()

	result := make(chan interface{})

	go func() {
		for {
			columnRanges := fmt.Sprintf("%s%d:%s%d", startChar, start, endChar, batchSize)
			sheet, err := svr.Spreadsheets.Values.Get(selectedSheet, columnRanges).Do()
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
	}()

	for e := range result {
		fmt.Fprint(w, e)
		if e == '\n' {
			w.Flush()
		}
	}

}

func getColumChar(columns int) string {
	return string('a' + (columns - 1))
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	scanCmd.Flags().Int("endColumn", 5, "Defines the End column")
	scanCmd.Flags().Int("startColumn", 1, "Define the start column")

	scanCmd.Flags().Int("start", 1, "Defines the starting row")
	scanCmd.Flags().Int("size", 10, "Define the BatchSize")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
