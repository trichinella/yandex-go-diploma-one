package erroring

import "errors"

var ErrIncorrectNumber = errors.New("incorrect number")
var ErrEmptyRequest = errors.New("empty request")
var ErrBadJSON = errors.New("bad json")
