package duobb

import (
	"errors"
)

var (
	ErrorNilData       error = errors.New("Nil data")
	ErrorIllegalData   error = errors.New("More than data len")
	ErrorDuobbMsg      error = errors.New("Error duobb msg")
	ErrorBase64Decode  error = errors.New("Base64 decode error")
	ErrorBase64Encode  error = errors.New("Base64 encode error")
	ErrorBadMethod     error = errors.New("Bad method")
	ErrorUnkownService error = errors.New("Unknown service")
)
