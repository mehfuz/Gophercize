package main

import (
	"fmt"

	"github.com/gophercises/gophercise17/crypton/cmd"
)

func main() {

	err := cmd.Rootcmd.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

}
