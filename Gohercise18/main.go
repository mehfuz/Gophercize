package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/gophercises/Gohercise18/primitive"
)

var GenImgFunc = GenImg
var TempFileFunc = ioutil.TempFile
var CopyFunc = io.Copy

// GenImg() calls the actual primitive routine to trasfom source image
// and returns the filename of the transformed file.
func GenImg(file io.Reader, ext string, numshapes int, mode primitive.Mode) (string, error) {
	output, err := primitive.TransformImg(file, ext, numshapes, primitive.ArgMode(mode))
	if err != nil {
		//http.Error(resp, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	saveout, err := getExtFile("", ext)
	if err != nil {
		//http.Error(resp, err.Error(), http.StatusInternalServerError)
		return "", err
	}
        _,err = CopyFunc(saveout, output)
        if err!=nil{
                return "",err
        }

	return saveout.Name(), err
}

// RenderSetChoice is executed if the url contains mode parameter
func RenderSetChoice(resp http.ResponseWriter, req *http.Request, ext string, mode primitive.Mode, fileSeeker io.ReadSeeker) {

	op := []OptStruct{
		{20, mode},
		{30, mode},
		{40, mode},
		{50, mode},
	}
	opFileList, err := GenImgList(ext, fileSeeker, op...)
	if err != nil {
                http.Error(resp, err.Error(), http.StatusInternalServerError)
                return
	}
	htmlist := `<html>
                        <body>
                {{range .}}
                <a href="/modify/{{.Name}}?mode={{.Mode}}&n={{.Numshapes}}">
                <img style ="width 30%" src="/pics/{{.Name}}">
                {{end}}
                </body>
                </html>
        `
	templ := template.Must(template.New("").Parse(htmlist))

	type Opts struct {
		Name      string
		Mode      primitive.Mode
		Numshapes int
	}
	var opts []Opts
	for index, val := range opFileList {
		opts = append(opts, Opts{Name: filepath.Base(val), Mode: op[index].mode, Numshapes: op[index].num})
	}

	// err = templ.Execute(resp, opts)
	// if err != nil {
	// 	panic(err)
	// }
                checkError(templ.Execute(resp,opts))
}
// for testing
func checkError(err error) {
        if err!=nil{
           fmt.Println(err)
           return
        }
}
// }
// RenderInitialChoices generates default choices of transformed images
// and pass them to responsewriter.
func RenderInitialChoices(resp http.ResponseWriter, req *http.Request, ext string, fileSeeker io.ReadSeeker) {
	op := []OptStruct{
		{22, primitive.TriangleMode},
		{22, primitive.CircleMode},
		{22, primitive.ComboMode},
		{22, primitive.PolygonMode},
	}
	opFileList, err := GenImgList(ext, fileSeeker, op...)
	if err != nil {
                http.Error(resp, err.Error(), http.StatusInternalServerError)
                return
	}
	htmlist := `<html>
                        <body>
                {{range .}}
                <a href="/modify/{{.Name}}?mode={{.Mode}}">
                <img style ="width 30%" src="/pics/{{.Name}}">
                {{end}}
                </body>
                </html>
        `
	templ := template.Must(template.New("").Parse(htmlist))
        
type Opts struct {
        Name string
        Mode primitive.Mode
}
	var opts []Opts
	for index, val := range opFileList {
		opts = append(opts, Opts{Name: filepath.Base(val), Mode: op[index].mode})
	}

	// err = templ.Execute(resp, opts)
	// if err != nil {
	// 	panic(err)
        // }
        checkError(templ.Execute(resp, opts))
        

}



type OptStruct struct {
	num  int
	mode primitive.Mode
}

// GenImgList generates list of images to be displayed
// when mode is not selected
func GenImgList(ext string, fileSeeker io.ReadSeeker, opts ...OptStruct) ([]string, error) {
	opFileList := []string{}
	for _, value := range opts {
		fileSeeker.Seek(0, 0)
		opFileName, err := GenImgFunc(fileSeeker, ext, value.num, value.mode)
		if err != nil {
			return nil, err
		}
		opFileList = append(opFileList, opFileName)
	}
	return opFileList, nil
}

func RootHandler(resp http.ResponseWriter, req *http.Request) {
	html := `<html>
                         <h3> UPLOAD IMAGE HERE</h3>
                         <form action="/upload" method="POST" enctype="multipart/form-data">
                         <input type="file" name ="img"/>
                         <button type="submit">UPLOAD</button>
                         </form>
                         </html>`

	fmt.Fprint(resp, html)
}
func UploadHandler(resp http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("img")
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)[1:]
	output, err := getExtFile("", ext)
	defer output.Close()
	if err != nil {
                http.Error(resp, err.Error(), http.StatusInternalServerError)
                return 
	}
	_, err = CopyFunc(output, file)
	if err != nil {
                http.Error(resp, err.Error(), http.StatusInternalServerError)
                return
	}

	http.Redirect(resp, req, "/modify/"+filepath.Base(output.Name()), http.StatusFound)
}
func ModifyHandler(resp http.ResponseWriter, req *http.Request) {
	file, err := os.Open("./pics/" + filepath.Base(req.URL.Path))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	rawmode := req.FormValue("mode")
	ext := filepath.Ext(file.Name())[1:]
	if rawmode == "" {
		//render initial choices
		RenderInitialChoices(resp, req, ext, file)
		return
	}
	// call to render selected choices
	mode, err := strconv.Atoi(rawmode)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
	}
	numStr := req.FormValue("n")
	if numStr == "" {
		RenderSetChoice(resp, req, ext, primitive.Mode(mode), file)
		return
	}
	_, err = strconv.Atoi(numStr) //temperorily
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "image/png")
	http.Redirect(resp, req, "/pics/"+filepath.Base(file.Name()), http.StatusFound)

}
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)

	mux.HandleFunc("/upload", UploadHandler)

	mux.HandleFunc("/modify/", ModifyHandler)

	//FileServer Configuration
	fileserver := http.FileServer(http.Dir("./pics"))
	mux.Handle("/pics/", http.StripPrefix("/pics", fileserver))
	fmt.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))

}

func getExtFile(prefix, suffix string) (*os.File, error) {
	infile, err := TempFileFunc("./pics/", prefix)
	if err != nil {
		return nil, err
	}
	defer os.Remove(infile.Name())
	fileName := fmt.Sprintf("%s.%s", infile.Name(), suffix)
	fmt.Println(fileName)
	return os.Create(fileName)
}
