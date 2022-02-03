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
	"github.com/lukas8219/gcsv/util"
	"github.com/spf13/cobra"
)

// favoriteCmd represents the favorite command
var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Add a Sheet to Favorite [link]",
	Long:  `Add a Sheet to Favorite. You can add it using the HTTPS or only the ID as Argument`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please select one of the options [add] [remove]")
	},
}

var favoriteAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a Sheet to Favorite [link]",
	Long:  `Add a Sheet to Favorite. You can add it using the HTTPS or only the ID as Argument`,
	Run:   add,
}

func add(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatalln("One and only one argument required [link]")
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Fatalln(err)
	}

	link := util.Parse(args[0])

	storage.Store(storage.FavoriteSheet{
		Name: name,
		ID:   link,
	})
}

func init() {
	rootCmd.AddCommand(favoriteCmd)

	favoriteCmd.AddCommand(favoriteAddCommand)

	favoriteAddCommand.Flags().String("name", "default", "Sets the name")
	favoriteAddCommand.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// favoriteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// favoriteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
