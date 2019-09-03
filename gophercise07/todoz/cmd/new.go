package cmd

import (
	"fmt"
	"strings"

	"github.com/gophercises/gophercise07/todoz/store"
	"github.com/spf13/cobra"
)


var Addnewtask = &cobra.Command{
	Use:   "new",
	Short: "Adds a new task in the list of our tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Provide task..")
		}
		task := strings.Join(args, " ")
		_, err := store.InsertTask(task)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("task \"%s\" added..:)\n", task)

	},
}

// this is executed before call to main
func init() {
	RootCmd.AddCommand(Addnewtask)
}
