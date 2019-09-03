package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"testing"
)

//Testing of the root command

func TestGet(test *testing.T) {
	args := []string{"testkey"}
	oldstdout := os.Stdout
	chanl := make(chan string)
	r, w, _ := os.Pipe()
	os.Stdout = w
	Getcmd.Run(Getcmd, args)

	go func() {
		var buf bytes.Buffer
		Setcmd.Run(Setcmd, args)
		io.Copy(&buf, r)
		chanl <- buf.String()
	}()

	w.Close()
	os.Stdout = oldstdout
	output := <-chanl
	fmt.Println(output)
	if !strings.Contains(output, "testvalue") {
		test.Error("Set cli test failed..")
	}
}

func TestErrGetFunc(test *testing.T) {
	Getcmd.Run(Getcmd, []string{"dumbkey"})
}
