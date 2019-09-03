package store

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	home "github.com/mitchellh/go-homedir"
)

//var testvar *testing.T

type teststruct struct {
	input    int
	expected []byte
}

// test boti() and itob()

var hdir, _ = home.Dir()

var path = filepath.Join(hdir, "todo.db")
var testval = teststruct{
	input:    5,
	expected: []byte{0, 0, 0, 0, 0, 0, 0, 5},
}

func TestItob(t *testing.T) {
	for index, value := range itob(5) {
		if value != testval.expected[index] {

			t.Error("Itob failed")

		}
	}
}

func TestBtoi(t *testing.T) {
	if btoi(testval.expected) != testval.input {
		t.Error("Btoi failed")
	}
}

func TestInsertTask(t *testing.T) {
	Init(path)
	_, err := InsertTask("Dummy task")
	if err != nil {
		t.Error("Insert Failed")
		fmt.Println(err)
	}
}

func TestRemoveTasks(t *testing.T) {
	err := RemoveTasks(5)
	if err != nil {
		log.Println("Delete failed")
	}
}

func TestGetAll(t *testing.T) {
	_, err := GetAll()
	if err != nil {
		log.Println("Get all failed")
	}
}

func TestInit(t *testing.T) {
	err := Init("/")
	if err != nil {
		log.Println("Failed init..")
	}
}

func TestErrInsertRecord(test *testing.T) {
	err := Init("test.db")
	DbCon.Close()
	_, err = InsertTask("test record..")
	if err != nil {
		fmt.Println("InsertRecord Error" + err.Error())
	}
}

