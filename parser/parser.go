package parser

import (
	"fmt"
	"strings"

	"github.com/kballard/go-shellquote"
	"keybot/api"
	"keybot/config"
)

type cmd struct {
	Command  string
	HelpText string
	CmdFunc  cmdFunc
	ShowHelp bool
	Active   bool
}

type cmdFunc func(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (CmdOut, error)

type CmdOut struct {
	Message string
	Channel api.Channel
}

var Commands = make(map[string]*cmd)

func RegisterCommand(command, helptext string, showhelp, active bool, cmdFunc cmdFunc) {
	Commands[command] = &cmd{
		Command:  command,
		HelpText: helptext,
		CmdFunc:  cmdFunc,
		ShowHelp: showhelp,
		Active:   active,
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
			args, err := GetArgs(msgText)
			if err != nil {
				fmt.Println(err)
				return
			}

			isActiveChannel := !chat.IsTeam
			_, isActiveTeam := c.ActiveTeams[chat.Name]
			if isActiveTeam && !isActiveChannel {
				_, isActiveChannel = c.ActiveTeams[chat.Name].ActiveChannels[chat.Channel]
			}
			cmd, isCommand := Commands[args[0]]
			if isCommand && isActiveChannel {
				m, err := cmd.CmdFunc(args, message, c)
				if err, ok := err.(*CmdError); ok {
					chat.SendMessage(fmt.Sprintf("%s: %s", err.Command, err.Message))
				} else {
					m.Channel.SendMessage(m.Message)
				}
			}
		}
	default:
		// do nothing
	}

}
