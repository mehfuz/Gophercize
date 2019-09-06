package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	//"github.com/alecthomas/chroma/quick"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var StyleValue = styles.Get("github")
var lex = lexers.Get("go")
var lexStringFunc = lex.Tokenise

func HandleDebug(resp http.ResponseWriter, req *http.Request) {
	filePath := req.FormValue("path")
	filenumber := req.FormValue("line")
	for index, val := range filenumber {
		if val == ' ' {
			filenumber = filenumber[:index]
		}
	}
	fmt.Println(filenumber)
	intnumber, err := strconv.Atoi(filenumber)
	if err != nil {
		fmt.Println(err.Error())
		intnumber = -1
		//return
	}
	data, err := os.Open(filePath)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
	bin := bytes.NewBuffer(nil)
	_, err = io.Copy(bin, data)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

	var lines [][2]int
	if intnumber > 0 {
		lines = append(lines, [2]int{intnumber, intnumber})
	}
	if StyleValue == nil {
		StyleValue = styles.Fallback
	}
	iter, err := lexStringFunc(nil, bin.String())
	formatter := html.New(html.TabWidth(2), html.WithLineNumbers(), html.HighlightLines(lines), html.LineNumbersInTable())
	resp.Header().Set("Content-Type", "text/html")
	formatter.Format(resp, StyleValue, iter)
	//err = quick.Highlight(resp, bin.String(), "go", "html", "monokai")
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}

}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, parseLines(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func Handlepanicmode(resp http.ResponseWriter, req *http.Request) {
	panic("Panic..")
}
func parseLines(response string) string {
	lines := strings.Split(response, "\n")

	for index, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		tline := ""
		for ind, ch := range line {
			if ch == ':' {

				tline = line[1:ind]
				break
			}
		}
		var strbuilder strings.Builder
		//add line number as parameter to the url
		for i := len(tline) + 2; i < len(line); i++ {
			if line[i] < '0' || line[i] > '9' {
				break
			}
			strbuilder.WriteByte(line[i])
		}

		V := url.Values{}
		V.Set("path", tline)
		V.Set("line", strbuilder.String())
		lines[index] = "\t<a href=\"/debug?" + V.Encode() + "\">" + tline + ":" + strbuilder.String() + "</a>" + line[len(tline)+2+len(strbuilder.String()):]
	}
	return strings.Join(lines, "\n")
}

var Mockhttp = http.ListenAndServe

func main() {
	Router := http.NewServeMux()
	Router.HandleFunc("/debug", HandleDebug)
	Router.HandleFunc("/panic", Handlepanicmode)
	fmt.Println("Server running on Port 3000")
	log.Fatal(Mockhttp(":3000", devMw(Router)))
}
