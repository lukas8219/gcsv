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

	"github.com/lukas8219/gcsv/storage"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Get a configuration property value",
	Long:  "Get a config pro",
	Run:   get,
}

func get(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("One and only one argument is required")
	}

	storage := storage.GetStorage()

	val, err := storage.GetProp(args[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Property %s Value is %s", args[0], val))
}

var setCommand = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration property",
	Long:  "Set a config prop",
	Run:   set,
}

func set(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("Two and only two arguments are required")
	}

	prop := args[0]
	value := args[1]

	storage := storage.GetStorage()

	err := storage.SetProp(prop, value)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setCommand)
	configCmd.AddCommand(getCommand)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
