package main
import (
	"testing"
	"errors"
)


func TestMainFunc(t *testing.T){
	handleError(errors.New("dummy error"))
	main()
	
}