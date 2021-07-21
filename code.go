package errors

import (
	"fmt"
	"git.cai-inc.com/support/errors/code"
	"github.com/novalagung/gubrak"
	"net/http"
	"sync"
)

var (
	unknownCoder ErrCode = ErrCode{1, http.StatusInternalServerError, "An internal server error occurred", "http://git.cai-inc.com/support/errors/README.md"}
)

// Coder defines an interface for an error code detail information.
type Coder interface {
	// HTTP status that should be used for the associated error code.
	HTTPStatus() int

	// External (user) facing error text.
	String() string

	// Reference returns the detail documents for user.
	Reference() string

	// Code returns the code of the coder
	Code() int
}


// codes contains a map of error codes to metadata.
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}


// MustRegister register a user define error code.
// It will panic when the same Code already exist.
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code '0' is reserved by 'github.com/marmotedu/errors' as ErrUnknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder parse any error into *withCode.
// nil error will return nil direct.
// None withStack error will be parsed as ErrUnknown.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode reports whether any error in err's chain contains the given error code.
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}


// ErrCode implements `git.cai-inc.com/support/errors`.Coder interface.
type ErrCode struct {
	// C refers to the code of the ErrCode.
	C int

	// HTTP status that should be used for the associated error code.
	HTTP int

	// External (user) facing error text.
	Ext string

	// Ref specify the reference document.
	Ref string
}

// Code returns the integer code of ErrCode.
func (coder ErrCode) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder ErrCode) String() string {
	return coder.Ext
}

// Reference returns the reference document.
func (coder ErrCode) Reference() string {
	return coder.Ref
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder ErrCode) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}
	return coder.HTTP
}
func init() {
	codes[unknownCoder.Code()] = unknownCoder
	register(code.ErrUserNotFound, 404, "User not found")
	register(code.ErrUserAlreadyExist, 400, "User already exist")
	register(code.ErrReachMaxCount, 400, "Secret reach the max count")
	register(code.ErrSecretNotFound, 404, "Secret not found")
	register(code.ErrSuccess, 200, "OK")
	register(code.ErrUnknown, 500, "Internal server error")
	register(code.ErrBind, 400, "Error occurred while binding the request body to the struct")
	register(code.ErrValidation, 400, "Validation failed")
	register(code.ErrTokenInvalid, 401, "Token invalid")
	register(code.ErrDatabase, 500, "Database error")
	register(code.ErrEncrypt, 401, "Error occurred while encrypting the user password")
	register(code.ErrSignatureInvalid, 401, "Signature is invalid")
	register(code.ErrExpired, 401, "Token expired")
	register(code.ErrInvalidAuthHeader, 401, "Invalid authorization header")
	register(code.ErrMissingHeader, 401, "The `Authorization` header was empty")
	register(code.ErrorExpired, 401, "Token expired")
	register(code.ErrPasswordIncorrect, 401, "Password was incorrect")
	register(code.ErrPermissionDenied, 403, "Permission denied")
	register(code.ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	register(code.ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	register(code.ErrInvalidJSON, 500, "Data is not valid JSON")
	register(code.ErrEncodingJSON, 500, "JSON data could not be encoded")
	register(code.ErrDecodingJSON, 500, "JSON data could not be decoded")
	register(code.ErrInvalidYaml, 500, "Data is not valid Yaml")
	register(code.ErrEncodingYaml, 500, "Yaml data could not be encoded")
	register(code.ErrDecodingYaml, 500, "Yaml data could not be decoded")
}
func register(code int, httpStatus int, message string, refs ...string) {
	found, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)
	if !found {
		panic("http code not in `200, 400, 401, 403, 404, 500`")
	}

	var reference string
	if len(refs) > 0 {
		reference = refs[0]
	}

	coder := &ErrCode{
		C:    code,
		HTTP: httpStatus,
		Ext:  message,
		Ref:  reference,
	}

	MustRegister(coder)
}
