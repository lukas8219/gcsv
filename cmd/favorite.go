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

// favoriteCmd represents the favorite command
var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Add a Sheet to Favorite [link]",
	Long:  `Add a Sheet to Favorite. You can add it using the HTTPS or only the ID as Argument`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Please select one of the options [add] [remove]")
	},
}

var favoriteListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all the Favorites",
	Long:  "List all the Favorites",
	Run:   list,
}

func list(cmd *cobra.Command, args []string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 8, '\t', 0)
	storage := storage.GetStorage()

	fmt.Fprintf(w, "Name\tID\n")

	values := storage.ListAll()

	for _, val := range values {
		fmt.Fprintf(w, "%s\t%s\n", val.Name, val.ID)
	}

	w.Flush()
}

var favoriteRemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a Sheet from Favorite",
	Long:  `Remove a Sheet from Favorite. Use with the name`,
	Run:   remove,
}

var favoriteAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a Sheet to Favorite [link]",
	Long:  `Add a Sheet to Favorite. You can add it using the HTTPS or only the ID as Argument`,
	Run:   add,
}

var favoriteSetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set a favorite",
	Long:  "Set a Favorite",
	Run:   setFavorite,
}

func setFavorite(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("One and only one argument is required [name]")
	}

	storage := storage.GetStorage()
	_, err := storage.Get(args[0])
	if err != nil {
		log.Fatal(err)
	}
	storage.SetSelectedFavorite(args[0])
	log.Printf("Selected Favorite Set to %s\n", args[0])
}

func remove(cmd *cobra.Command, args []string) {
	storage := storage.GetStorage()
	if len(args) != 1 {
		log.Fatal("One and only one argument required")
	}

	id := args[0]

	_, err := storage.Get(id)
	if err != nil {
		log.Fatalln("There's no Sheet with this Name")
	}

	err = storage.Remove(id)
	if err != nil {
		log.Fatal(err)
	}
}

func add(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatalln("One and only one argument required [link]")
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Fatalln(err)
	}
	database := storage.GetStorage()

	id, err := database.Get(args[0])
	if err != nil {
		id = util.ParseLink(args[0])
	}

	database.Store(storage.FavoriteSheet{
		Name: name,
		ID:   id,
	})
}

func init() {
	rootCmd.AddCommand(favoriteCmd)

	favoriteCmd.AddCommand(favoriteAddCommand)
	favoriteCmd.AddCommand(favoriteRemoveCmd)
	favoriteCmd.AddCommand(favoriteListCommand)
	favoriteCmd.AddCommand(favoriteSetCommand)

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
