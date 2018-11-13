package parser

import (
	"strings"
	
	"github.com/kballard/go-shellquote"
)

func GetArgs(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	cmdArgs, err := shellquote.Split(s)
	if err != nil {
		return nil, err
	}
	if len(cmdArgs) == 0 {
		return nil, &cmdMissingError{"No command provided"}
	}
	return cmdArgs, nil
}
