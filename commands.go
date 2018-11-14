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
func cmdHelp(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	for _, command := range parser.Commands {
		if command.ShowHelp {
			response += fmt.Sprintf("`%s%s`\n  %s\n", CommandPrefix, command.Command, command.HelpText)
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

func cmdPing(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
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

func cmdConfig(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)

	if message.Msg.Sender.Username != config.BotOwner {
		return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You do not have permission to configure this bot.", message.Msg.Sender.Username)}
	}

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
			return parser.CmdOut{}, &cmdError{args[0], "'botOwner' must be set in config file directly."}
		default:
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - No such variable.", args[2])}
		}
	case "get":
		if len(args) < 3 {
			return parser.CmdOut{}, &cmdError{args[0], "Missing arguments."}
		}
		switch strings.ToLower(args[2]) {
		case "botowner":
			response = fmt.Sprintf("botOwner: %s", config.BotOwner)
		default:
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - No such variable.", args[2])}
		}
	case "blacklist":
		switch strings.ToLower(args[2]) {
		case "add":
			if len(args) < 4 {
				return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - No users provided.", args[1])}
			}
			if config.Blacklist == nil {
				config.Blacklist = make(map[string]struct{})
				config.Write()
			}

			var (
				added string
				notadded string
			)
		case "remove":
		case "read":
			if config.Blacklist == nil {
				config.Blacklist = make(map[string]struct{})
				config.Write()
			}
			if len(config.Blacklist) > 0 {
				response = "Blacklisted users:\n```\n"
				for user, _ := range config.Blacklist {
					response += fmt.Sprintln(user)
				}
				response += "```"
			} else {
				response = "There are no blacklisted users."
			}
		default:
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - Invalid action.", args[3])}
		}
	default:
		return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - Invalid command.", args[1])}
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

func commandHandler(message api.ChatMessageIn, c *config.ConfigJSON) {
	// Get channel details
	var chat = api.Channel{Name: message.Msg.Channel.Name}
	if message.Msg.Channel.MembersType == "team" {
		chat.IsTeam = true
		chat.Channel = message.Msg.Channel.TopicName
	}

	switch strings.ToLower(message.Msg.Content.MsgType) {
	case "text":
		msgText := strings.TrimSpace(message.Msg.Content.Text.Body)
		if strings.HasPrefix(msgText, CommandPrefix) {
			msgText = strings.TrimPrefix(msgText, CommandPrefix)
			args, err := parser.GetArgs(msgText)
			if err != nil {
				fmt.Println(err)
				return
			}

			isActiveChannel := !chat.IsTeam
			_, isActiveTeam := c.ActiveTeams[chat.Name]
			if isActiveTeam && !isActiveChannel {
				_, isActiveChannel = c.ActiveTeams[chat.Name].ActiveChannels[chat.Channel]
			}
			cmd, isCommand := parser.Commands[args[0]]
			if isCommand && isActiveChannel {
				m, err := cmd.CmdFunc(args, message, c)
				if err, ok := err.(*cmdError); ok {
					chat.SendMessage(fmt.Sprintf("%s: %s", err.command, err.message))
				} else {
					m.Channel.SendMessage(m.Message)
				}
			}
		}
	default:
		// do nothing
	}

}
