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

func TestSet(test *testing.T) {
	args := []string{"testkey", "testvalue"}
	oldstdout := os.Stdout
	chanl := make(chan string)
	r, w, _ := os.Pipe()
	os.Stdout = w
	Setcmd.Run(Setcmd, args)

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
	if !strings.Contains(output, "Data successfully encoded.") {
		test.Error("Set cli test failed..")
	}
}

func TestErrSetfunc(test *testing.T) {

	Setcmd.Run(Setcmd, []string{"", ""})
	filename = "/"
	Setcmd.Run(Setcmd, []string{"c", "k"})

}

//var s string
// b := make([]byte, 4)
// //buf := bytes.NewBuffer(b)

// Setcmd.Run(Setcmd, args)

// //result, err := ioutil.ReadAll(os.Stdout.
// _, err := os.Stdout.Read(b)

// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println("------------->", string(b))
//  var buf bytes.Buffer
// Rootcmd.AddCommand(Setcmd)
//
// Rootcmd.SetOut(&buf)
// Rootcmd.SetArgs(args)
// _,err := Rootcmd.ExecuteC()
// if err!=nil{
// 	fmt.Println(err)
// }
//   Rootcmd.ExecuteC()
// fmt.Println("------->>>>"+buf.String())
// cmd := exec.Command("crypton","set","testkey","testvalue")
// stdout,err :=cmd.StdoutPipe()
// if err !=nil{
// 	test.Error(err)
// }
// if err := cmd.Start(); err != nil {
// 	fmt.Println(err)
// }
// io.Copy(&buf,stdout)
// cmd.Wait()
// fmt.Println("------>>>>>"+ buf.String())
