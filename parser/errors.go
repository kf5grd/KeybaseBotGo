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

type CmdError struct {
	Command string
	Message string
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("%s: %s", e.Command, e.Message)
}

type cmdMissingError struct {
	message string
}

func (e *cmdMissingError) Error() string {
	return fmt.Sprint(e.message)
}
