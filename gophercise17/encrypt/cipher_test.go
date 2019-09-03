package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

var key = "This is key"
var text = "This is my secret test"
var prevEnc = "2ef309bb9f2c2f06a46b109858872df4ef38a199a9e2b57d6e23f191cb315499ba1ac19857fa"

//

func TestEncDec(test *testing.T) {
	encText, err := Enc(key, text)

	if err != nil {
		test.Error(err)
	}

	var plaintext string
	plaintext, err = Dec(key, encText)

	if err != nil {
		test.Error(err)
	}
	//fmt.Println(plaintext)
	if plaintext != text {
		fmt.Printf("Failed.Expeced %s got %s", text, plaintext)
		//errors.New(fmt.printf("Failed.Expeced %s got %s", text, plaintext))
	}
}

func TestEncWriterErr(test *testing.T) {

	oldGetencStream := GetEncryptStreamFunc
	GetEncryptStreamFunc = func(key string, iv []byte) (cipher.Stream, error) {
		return nil, errors.New("Thrown error")
	}
	F, _ := os.Open("testing.txt")
	_, err := EncWriter("sme-key", F) //////////
	if err != nil {
		log.Println(err)
	}
	GetEncryptStreamFunc = oldGetencStream

	oldIoReadFullFunc := IoReadFullFunc
	IoReadFullFunc = func(r io.Reader, buf []byte) (n int, err error) {
		return 0, errors.New("io.read mocked")
	}
	_, err = EncWriter("sme-key", F) //////////
	if err != nil {
		log.Println(err)
	}
	_,err = Enc("sme-key","smeval")
	IoReadFullFunc = oldIoReadFullFunc

	oldGetDecryptStreamFunc := GetDecryptStreamFunc

	GetDecryptStreamFunc = func(key string, iv []byte) (cipher.Stream, error) {
		return nil, errors.New("Decrypt function mocked")
	}
	_, err = DecReader("smekey", F)
	if err != nil {
		log.Println(err)
	}
	GetDecryptStreamFunc = oldGetDecryptStreamFunc
	F.Close()
	_, err = EncWriter("sme-key", F) //////////
	if err != nil {
		log.Println(err)
	}
	_, err = DecReader("smekey", F)
	if err != nil {
		log.Println(err)
	}

}

func TestGetEncryptDecryptStreamErr(test *testing.T) {
	oldgetCipherBlock := getCipherBlockFunc
	getCipherBlockFunc = func(key string) (cipher.Block, error) {
		return nil, errors.New("mocker cipherblock")
	}
	iv := make([]byte, aes.BlockSize)
	_, err := GetEncryptStream("vault", iv)
	if err != nil {
		log.Println(err)
	}
	_, err = GetDecryptStream("valut", iv)
	if err != nil {
		log.Println(err)
	}
	_, err = Enc("Smekey", "smeval")
	if err != nil {
		log.Println(err)
	}
	_, err = Dec("Smekey", "smeval")
	if err != nil {
		log.Println(err)
	}
	getCipherBlockFunc = oldgetCipherBlock
}

func TestMockedDecodeFunc(test *testing.T) {
	_, err := Dec("", "1")
	if err != nil {
		log.Println(err)
	}
	et, err := Enc("Smekey", "smeval")
	if err != nil {
		log.Println(err)
	}
	oldAesBlockSize := AesBlockSize
	defer func() { AesBlockSize = oldAesBlockSize }()
	AesBlockSize = 400
	ot, err := Dec("Smekey", et)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("original text:::" + ot)
}


func TestMockedRandFunc(test *testing.T){

}
