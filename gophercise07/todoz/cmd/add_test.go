package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gophercises/gophercise07/todoz/store"
	home "github.com/mitchellh/go-homedir"
)

var hdir, _ = home.Dir()

var path = filepath.Join(hdir, "todo.db")

func TestAdd(t *testing.T) {
    store.Init(path)
	args := []string{"Add", "New", "value"}
    a := []string{}
	Addnewtask.Run(Addnewtask, args)
	store.DbCon.Close()
	//store.Init("/")
	Addnewtask.Run(Addnewtask, a)
}
