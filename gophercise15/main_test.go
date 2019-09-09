package main

import (
	"runtime/debug"
	"testing"
	"time"

	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/alecthomas/chroma"

	"strings"
)

type fn func(resp http.ResponseWriter, rq *http.Request)

//Test Function 
func CheckLinks(endpoint fn, method string, url string, query string, expectedStatus int) string {
	req, err := http.NewRequest(method, url+query, nil)
	if err != nil {
		fmt.Println(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc(url, endpoint)
	handler := devMw(mux)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		fmt.Printf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	return rr.Body.String()
}

func TestHandleDebug(t *testing.T) {
	query := "?path=/home/mehfuz/go/src/github.com/Gophercize-master/gophercise15/main.go"
	responsestirng := CheckLinks(HandleDebug, "GET", "/debug", query, 200)
	//fmt.Println(responsestirng)
	// Check the response body is what we expect.
	b, err := ioutil.ReadFile("/home/mehfuz/go/src/github.com/Gophercize-master/gophercise15/main.go")
	if err != nil {
		t.Error("Reading expected file error.")
	}

	expected := string(b)

	if strings.Contains(responsestirng, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responsestirng, expected)
	} else {
		fmt.Println("Test success...")
	}
	//main()
} //Just twik the file name in next case.

func TestPanic(test *testing.T) {

	_ = CheckLinks(Handlepanicmode, "GET", "/panic", "", 500)
	//fmt.Println("_________" + resp)
}

func TestPostpanicLink(test *testing.T) {
	RespString := CheckLinks(HandleDebug, "GET", "/debug", "?line=24+%2B0xa4&path=c%3A%2Fgo%2Fsrc%2Fruntime%2Fdebug%2Fstack.go", 200)
	if strings.Contains(RespString, "func PrintStack()") {
		test.Error("Test Failed" + fmt.Sprintf("expected func PrintStack(), got %s", RespString))
	}
}

func TestInvalidfileHandling(test *testing.T) {
	response := CheckLinks(HandleDebug, "GET", "/debug", "?path=cK", 500)
	fmt.Println("----->" + response)
}
func TestParseLines(t *testing.T) {
	str := parseLines("Test case coverage\n : \tcomplete")
	if str != "Test case coverage\n : \tcomplete" {
		t.Error("Test failed.......")
	}
}

func TestStylesGet(test *testing.T) {
	oldvalue := StyleValue
	defer func() { StyleValue = oldvalue }()
	StyleValue = nil
	_ = CheckLinks(HandleDebug, "GET", "/debug", "?path=cK", 500)
}

func TestNegativeLex(test *testing.T) {
	oldlexStringFunc := lexStringFunc
	defer func() {
		lexStringFunc = oldlexStringFunc
	}()
	lexStringFunc = func(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
		return nil, errors.New("mocked Lexer")
	}
	_ = CheckLinks(HandleDebug, "GET", "/debug", "?path=/home/mehfuz/go/src/github.com/Gophercize-master/gophercise15/main.go", 500)
}

func TestCreateLinks(t *testing.T) {
	stack := debug.Stack()
	link := parseLines(string(stack))
	if link == "" {
		t.Error("Expected link got", link)
	}
}

func TestMainFunc(test *testing.T) {
	go main()
	time.Sleep(1 * time.Second)
}
