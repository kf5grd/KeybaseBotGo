package parser

import (
	"strings"
	
	"keybot/api"
	"keybot/config"
	"github.com/kballard/go-shellquote"
)

type cmd struct {
	Command   string
	HelpText  string
	CmdFunc   cmdFunc
	ShowHelp  bool
	Active    bool
}

type cmdFunc func(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (CmdOut, error)

type CmdOut struct {
	Message string
	Channel api.Channel
}

var Commands = make(map[string]*cmd)

func RegisterCommand(command, helptext string, showhelp, active bool, cmdFunc cmdFunc) {
	Commands[command] = &cmd{
		Command:   command,
		HelpText:  helptext,
		CmdFunc:   cmdFunc,
		ShowHelp:  showhelp,
		Active:    active,
	}
}

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
