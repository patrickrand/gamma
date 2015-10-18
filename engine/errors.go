package engine

import (
	"bytes"
)

type Errors []error

func NewErrors(err ...error) Errors {
	return Errors(err)
}

func (errs Errors) Error() string {
	var buffer bytes.Buffer
	for e := range errs {
		buffer.WriteString(errs[e].Error())
	}
	return buffer.String()
}
