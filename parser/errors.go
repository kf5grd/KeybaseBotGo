package parser

import (
	"fmt"
)

type argError struct {
	message string
}

func (e *argError) Error() string {
	return fmt.Sprint(e.message)
}

type cmdMissingError struct {
	message string
}

func (e *cmdMissingError) Error() string {
	return fmt.Sprint(e.message)
}
