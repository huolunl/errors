package errors

import (
	"fmt"
	"git.cai-inc.com/support/errors/code"
	"testing"
)

func TestWithCode(t *testing.T) {
	err := WithCode(code.ErrDecodingJSON,"abc")
	coder := ParseCoder(err)
	fmt.Println(coder.String())
	fmt.Println(coder.Code())
	fmt.Println(coder.HTTPStatus())
}
