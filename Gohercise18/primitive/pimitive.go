package primitive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	//"log"
	"os/exec"
	"strings"
)

type Mode int

const (
	ComboMode Mode = iota
	TriangleMode
	RectangleMode
	EllipseMode
	CircleMode
	RotatedrectMode
	BeziersMode
	RotatedellipseMode
	PolygonMode
)

var TempFileFunc = ioutil.TempFile
// ArgMode is used set the mode for building image
func ArgMode(mode Mode) func() []string {

	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func getExtFile(prefix, suffix string) (*os.File, error) {
	infile, err := TempFileFunc("", prefix)
	if err != nil {
		return nil, err
	}
	defer os.Remove(infile.Name())
	fileName := fmt.Sprintf("%s.%s", infile.Name(), suffix)
	return os.Create(fileName)
}

var CopyFunc = io.Copy
var getExtFileFunc = getExtFile

// TransformImg converts the image from imgrdr using primitive and
// returns a reader to formed image..
func TransformImg(imgreadr io.Reader, ext string, numOfShapes int, options ...func() []string) (io.Reader, error) {
	var args []string
	for _, option := range options {
		args = append(args, option()...)
	}
	transformreader := bytes.NewBuffer(nil) // return value..
	infile, err := getExtFile("inp_", ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(infile.Name())

	opfile, err := getExtFileFunc("op_", ext)
	if err != nil {
		return nil, err
	}

	defer os.Remove(opfile.Name())

	_, err = io.Copy(infile, imgreadr)
	if err != nil {
		return nil, err //fmt.Println(err)
	}
	op, err := primitive(infile.Name(), opfile.Name(), numOfShapes, args...)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(op))
	_, err = CopyFunc(opfile, transformreader)
	if err != nil {
		return nil, err
	}
	return opfile, nil
}

// primitive function creates a routine for building image from primitive
// by embedding the parameters .
func primitive(input, output string, numShapes int, args ...string) ([]byte, error) {

	CmdString := fmt.Sprintf("-i %s -o %s -n %d ", input, output, numShapes) //"primitive", strings.Fields("-i samurai.png -o sam.png -n 100 -m 8")
	args = append(strings.Fields(CmdString), args...)
	cmd := exec.Command("primitive", args...)
	return cmd.CombinedOutput()
}
