package primitive

import (
	"io"
	"log"
	"os"
	"errors"
	"testing"
)

// func TestArgMode(test *testing.T) {
// 	expected := []string{"-m", "1"}
// 	fnc := ArgMode(TriangleMode)
// 	output := fnc()
// 	if output[0] != expected[0] && output[1] != expected[1] {
// 		test.Error("Test for Argmode failed")
// 	}
// }

func ChkFile(filename string) bool {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false //fmt.Println("file does not exist")
	}
	return true
}

// func TestPrimitive(test *testing.T) {
// 	_, err := primitive("test_samurai.png", "out.png", 5)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	//details, := os.Stat("./out.png")
// 	if !ChkFile("./out.png")  {
// 		test.Error("file not created..primitive failed")
// 	}
// }

func TestTransformImg(test *testing.T) {
	img_file, err := os.Open("test_samurai.png")
	if err != nil {
		log.Println(err)
		
	}
	output, err := TransformImg(img_file, "png", 4, ArgMode(EllipseMode))
	if err != nil {
		log.Println(err)
		
	}
	of, err := os.OpenFile("otrans.png", os.O_RDWR|os.O_CREATE, 0755)
	defer of.Close()
	io.Copy(of, output)
	details,err := os.Stat("otrans.png")
	if err != nil {
		log.Println(err)
		
	}
	if !ChkFile("otrans.png") && details.Size()>0{
		test.Error("Tansform function Failed")
	}

	output, err = TransformImg(img_file, "txt", 4, ArgMode(EllipseMode))
	if err != nil {
		log.Println(err)
		
	}

}


func TestMockedCopyFunc(test *testing.T){
	f,_ := os.Open("dummy.txt")
	_,err := TransformImg(f,"png",2,ArgMode(TriangleMode))
	if err!=nil{
		log.Println(err)
	 }
	oldCopyFunc :=CopyFunc
	defer func(){CopyFunc = oldCopyFunc}()
	CopyFunc = func (w io.Writer,r io.Reader)(int64 , error){
		return 0,errors.New("Mocked Copy Function.")
	}
	f,_ = os.Open("test_samurai.png")
	_,err = TransformImg(f,"png",2,ArgMode(TriangleMode))
	if err!=nil{
		log.Println(err)
	}

}


func TestMockgetExtFile(test *testing.T){
	f,_ := os.Open("test_samurai.png")
	oldgetExtFileFunc := getExtFileFunc
	defer func(){getExtFileFunc = oldgetExtFileFunc}()
	getExtFileFunc = func(a,b string)(*os.File,error){
		return nil,errors.New("Mocked getExtFileFunc()")
	}
	_,err := TransformImg(f,"png",2,ArgMode(TriangleMode))
	if err!=nil{
		log.Println(err)
	}

}
func TestTempFileMock(test *testing.T){
	f,_ := os.Open("test_samurai.png")
	oldTempFileFunc :=TempFileFunc
	defer func(){TempFileFunc = oldTempFileFunc}()
	TempFileFunc = func (a,b string) (*os.File,error){
		return nil,errors.New("Mocked ioutils.Tempfile")
	}
	_,err := getExtFile("a","png")
	if err!=nil{
		log.Println(err)
	}
	_,err = TransformImg(f,"png",2,ArgMode(TriangleMode))
	if err!=nil{
		log.Println(err)
	}
}