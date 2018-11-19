package main

import (
	"fmt"
	
	"keybot/api"
	"keybot/config"
	"keybot/parser"
)

func cmdHelp(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	for _, command := range parser.Commands {
		if command.ShowHelp {
			response += fmt.Sprintf("`%s%s`\n  %s\n", config.CommandPrefix, command.Command, command.HelpText)
		}
	}

	switch message.Msg.Channel.MembersType {
	case "team":
		channel = api.Channel{true, message.Msg.Channel.Name, message.Msg.Channel.TopicName}
	default:
		channel = api.Channel{Name: message.Msg.Channel.Name}
	}
	return parser.CmdOut{response, channel}, nil
}

func init() {
	parser.RegisterCommand("help", "", false, true, cmdHelp)
}
