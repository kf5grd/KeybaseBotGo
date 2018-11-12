package parser

import (
	"strings"
	
	"github.com/kballard/go-shellquote"
)

func GetArgs(s string) []string {
	s = strings.TrimSpace(s)
	cmdArgs, _ := shellquote.Split(s)
	return cmdArgs
}
