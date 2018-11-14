package main

import (
	"keybot/api"
	"keybot/config"
	"keybot/parser"

	"fmt"
	"strings"
)

type cmdError struct {
	command string
	message string
}

func (e *cmdError) Error() string {
	return fmt.Sprintf("%s: %s", e.command, e.message)
}
func cmdHelp(args []string, message api.ChatMessageIn, config config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	for _, command := range parser.Commands {
		if command.ShowHelp {
			response += fmt.Sprintf("`%s%s` - %s\n", CommandPrefix, command.Command, command.HelpText)
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

func cmdPing(args []string, message api.ChatMessageIn, config config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	if len(args) == 1 {
		response = "pong"
	} else {
		response = "ping "
		for _, u := range args[1:] {
			response += fmt.Sprintf("@%s, ", u)
		}
		response = strings.TrimSuffix(response, ", ")
	}
	
	switch message.Msg.Channel.MembersType {
	case "team":
		channel = api.Channel{true, message.Msg.Channel.Name, message.Msg.Channel.TopicName}
	default:
		channel = api.Channel{Name: message.Msg.Channel.Name}
	}
	return parser.CmdOut{response, channel}, nil
}

func cmdConfig(args []string, message api.ChatMessageIn, config config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	if len(args) < 2 {
		return parser.CmdOut{}, &cmdError{args[0], "Missing arguments."}
	}

	switch strings.ToLower(args[1]) {
	case "set":
		if len(args) < 3 {
			return parser.CmdOut{}, &cmdError{args[0], "Missing arguments."}
		}
		switch strings.ToLower(args[2]) {
		case "botowner":
		}
	case "get":
		if len(args) < 3 {
			return parser.CmdOut{}, &cmdError{args[0], "Missing arguments."}
		}
		switch strings.ToLower(args[2]) {
		case "botowner":
			response = fmt.Sprintf("botOwner: %s", config.BotOwner)
		}
	default:
		return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - Invalid command.")}
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
	parser.RegisterCommand("ping", "Responds with 'pong'", true, true, cmdPing)
	parser.RegisterCommand("config", "Get and set config values.", true, true, cmdConfig)
}
