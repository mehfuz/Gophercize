package cmd

import (
	"fmt"
	//"os"

	"github.com/gophercises/gophercise07/todoz/store"
	"github.com/spf13/cobra"
)

var ShowAll = store.GetAll
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Gives the list of all the pending tasks..",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("list called")
		List, err := ShowAll()
		if err != nil {
			fmt.Println(err)
			return //os.Exit(1)
		}
		if len(List) == 0 {
			fmt.Println("You have no pending tasks..")
		}
		for id, task := range List {
			fmt.Printf("%d . %s (%d)\n", id+1, task.Task, task.Id)
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
