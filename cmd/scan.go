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
	client := getClient()
	ctx := context.Background()

	svr, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)
	}

	//Try implementing with GoRoutines
	//Do an BinarySearch-like algorithm to find all Cells
	//Request a Batch from a fixed size column
	//Cache it
	//Duplicate the batch size, starting from the previous position
	//Repeat until no more batches return

	sheet, err := svr.Spreadsheets.Values.BatchGet("1HksEbnev-LX3T5gUfFR_2c9HLL8lodyMRv5q_CuYdnw").Ranges("A1:F10").Do()
	if err != nil {
		log.Fatalln(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)

	for _, val := range sheet.ValueRanges {
		for _, rowVal := range val.Values {
			for _, val := range rowVal {
				fmt.Fprint(w, val)
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, "\n")
		}
	}
	w.Flush()
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
