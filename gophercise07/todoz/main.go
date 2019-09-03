package main

import (
	"fmt"

	"github.com/gophercises/gophercise07/todoz/cmd"
	"github.com/gophercises/gophercise07/todoz/store"
)

func main() {

	path := "todo.db"

	handleError(store.Init(path))

	handleError(cmd.RootCmd.Execute())
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
