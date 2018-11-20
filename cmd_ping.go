package main

import (
	"keybot/api"
	"keybot/config"
	"keybot/parser"
)

func cmdPing(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel  api.Channel
		response string
	)

	response = "pong"

	switch message.Msg.Channel.MembersType {
	case "team":
		channel = api.Channel{true, message.Msg.Channel.Name, message.Msg.Channel.TopicName}
	default:
		channel = api.Channel{Name: message.Msg.Channel.Name}
	}
	return parser.CmdOut{response, channel}, nil
}

func init() {
	parser.RegisterCommand("ping", "Responds with 'pong'", true, true, cmdPing)
}
