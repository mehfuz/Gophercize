package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gophercises/Gohercise18/primitive"
)

type fn func(resp http.ResponseWriter, req *http.Request)

var File, _ = os.Open("./samurai.png")

func CheckLinks(endpoint fn, method string, url string, query string, expectedStatus int, hasbody bool, key string) string {
	var req *http.Request
	var err error
	if hasbody {
		file := File
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(key, "dummy.png")
		if err != nil {
			log.Println(err)
		}
		io.Copy(part, file)
		writer.Close()
		req, err = http.NewRequest(method, url+query, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		req, err = http.NewRequest(method, url+query, nil)
	}
	if err != nil {
		fmt.Println(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(endpoint)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedStatus {
		fmt.Printf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}
	return rr.Body.String()
}

func TestRoot(test *testing.T) {
	CheckLinks(RootHandler, "GET", "/", "", 200, false, "")
}

func TestUpoadLink(test *testing.T) {
	//check for folder before
	CheckLinks(UploadHandler, "POST", "/upload", "", 200, true, "img")
	//Check for folder after
	CheckLinks(UploadHandler, "POST", "/upload", "", 200, false, "")
}
func TestModifyLink(test *testing.T) {
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png", "", 200, true, "img")
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png?mode=2", "", 200, true, "img")
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png?mode=3&n=2", "", 200, true, "img")
}

func TestModifyNegative(test *testing.T) {
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png?mode=x", "", 500, true, "img")
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png?mode=2&n=u", "", 500, true, "img")

	oldGenImgFunc := GenImgFunc
	defer func() { GenImgFunc = oldGenImgFunc }()
	GenImgFunc = func(file io.Reader, ext string, numshapes int, mode primitive.Mode) (string, error) {
		return "", errors.New("Mocker GenImg")
	}
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png", "", 500, true, "img")
	CheckLinks(ModifyHandler, "GET", "/modify/155099759.png?mode=2", "", 500, true, "img")

}

func TestImproperFile(test *testing.T) {
	CheckLinks(ModifyHandler, "GET", "/modify/nodata.png?mode=2", "", 500, true, "img")
}
func TestGetExtfileErr(test *testing.T) {
	oldTempFileFunc := TempFileFunc

	TempFileFunc = func(a, b string) (*os.File, error) {
		return nil, errors.New("Mocked ExtFile")

	}
	getExtFile("", "")
	f, _ := os.Open("samurai.png")
	_, err := GenImg(f, "png", 5, primitive.TriangleMode)
	if err != nil {
		log.Println(err)
	}
	CheckLinks(UploadHandler, "POST", "/upload", "", 200, true, "img")
	TempFileFunc = oldTempFileFunc
	
	CheckLinks(UploadHandler, "POST", "/upload", "", 500, true, "test")

	_,err = GenImg(File,"txt",3,primitive.ComboMode)
	if err!=nil{
		log.Println(err)
	}


}

func TestErrCopy(test *testing.T){
	oldCopyFunc := CopyFunc
	defer func(){
		CopyFunc = oldCopyFunc
	}()
	CopyFunc = func (w io.Writer,r io.Reader)(int64 , error){
		return 0,errors.New("Mocked Copy Function.")
	}
	f, _ := os.Open("samurai.png")
	_, err := GenImg(f, "png", 5, primitive.TriangleMode)
	if err != nil {
		log.Println(err)
	}

	CheckLinks(UploadHandler, "POST", "/upload", "", 200, true, "img")
	checkError(errors.New("mocked it.."))

}

func TestMainFunc(test *testing.T) {
	go main()
	time.Sleep(1 * time.Second)
}
