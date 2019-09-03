package encrypt

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"log"
	"testing"
)

//mocked functions
type testStruct struct {
	key, val string
}

var testcase = []testStruct{
	{"key1", "Dummy string1"},
	{"key2", "dummy string 2"},
	{"key3", "Tesitng vaoflasmfms"},
}
var dummyvault = NewVault("It says this is not supposed to be smallaaaaaaaaa", "testing.txt")

func TestSetAndGet(test *testing.T) {
	for _, values := range testcase {
		err := dummyvault.Set(values.key, values.val)
		if err != nil {
			test.Error(err)
		}

	}
	for _, values := range testcase {
		val, err := dummyvault.Get(values.key)
		if err != nil {
			test.Error(err)
			break
		}
		if val != values.val {
			fmt.Printf("Expected %s and got %s \n", values.val, val)
			test.Error("Test failed ")
		}
	}

}

func TestNegativeSetAndGet(test *testing.T) {
	dumbvault := NewVault("dumbk", "/")
	err := dumbvault.Set("k", "val")
	if err != nil {
		log.Println(err.Error()) //test.Error(err)
	}
	_, err = dumbvault.Get("k")
	if err != nil {
		log.Println(err.Error())
	}
	oldDecReader := DecReaderFunc // Testing for failed
	DecReaderFunc = func(s string, r io.Reader) (*cipher.StreamReader, error) {
		return nil, errors.New("mocked DecReader")
	}

	err = dummyvault.Set("dumbkey", "dumbval")
	if err != nil {
		log.Println(err.Error()) //test.Error(err)
	}
	_, err = dummyvault.Get("dumbkey")
	if err != nil {
		log.Println(err.Error())
	}
	DecReaderFunc = oldDecReader

	///////////////////////////////////////////////////////////////////////////////////////
	oldEncwriter := EncWriterFunc
	EncWriterFunc = func(k string, w io.Writer) (*cipher.StreamWriter, error) {
		return nil, errors.New("mocker EncWriter")
	}
	err = dummyvault.Set("k", "val")
	if err != nil {
		log.Println(err.Error())
	}
	EncWriterFunc = oldEncwriter
}
