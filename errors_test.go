package errors

import (
	"fmt"
	"github.com/huolunl/errors/code"
	"testing"
)

func TestWithCode(t *testing.T) {
	err := WithCode(code.ErrDecodingJSON, "abc")
	coder := ParseCoder(err)
	fmt.Println(coder.String())
	fmt.Println(coder.Code())
	fmt.Println(coder.HTTPStatus())
	fmt.Println(coder.Reference())

	fmt.Println("---------------")
	// (# json), (+ detail),()
	fmt.Println(fmt.Sprintf("%+v", err))
	fmt.Println(fmt.Sprintf("%#v", err))
	fmt.Println(fmt.Sprintf("%#+v", err))
}

func TestRegister(t *testing.T) {
	Register(code.ErrDatabase, 500, "abc")
}
