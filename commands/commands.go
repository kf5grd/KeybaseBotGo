package commands

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

func CommandHandler(message api.ChatMessageIn, c *config.ConfigJSON) {
	// if user is blacklisted, do nothing
	if _, ok := c.Blacklist[message.Msg.Sender.Username]; ok {
		return
	}

	// Get channel details
	var chat = api.Channel{Name: message.Msg.Channel.Name}
	if message.Msg.Channel.MembersType == "team" {
		chat.IsTeam = true
		chat.Channel = message.Msg.Channel.TopicName
	}

	switch strings.ToLower(message.Msg.Content.MsgType) {
	case "text":
		msgText := strings.TrimSpace(message.Msg.Content.Text.Body)
		if strings.HasPrefix(msgText, c.CommandPrefix) {
			msgText = strings.TrimPrefix(msgText, c.CommandPrefix)
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
