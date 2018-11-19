package commands

import (
	"fmt"
	"strings"

	"keybot/api"
	"keybot/config"
	"keybot/parser"
)

func cmdUser(args []string, message api.ChatMessageIn, config *config.ConfigJSON) (parser.CmdOut, error) {
	var (
		channel api.Channel
		response string
	)
	switch message.Msg.Channel.MembersType {
	case "team":
		channel = api.Channel{true, message.Msg.Channel.Name, message.Msg.Channel.TopicName}
	default:
		channel = api.Channel{Name: message.Msg.Channel.Name}
		return parser.CmdOut{}, &cmdError{args[0], "Command can only be called from inside team."}
	}

	t := message.Msg.Channel.Name
	u := message.Msg.Sender.Username
	privs := make(map[string]bool)

	if (u == config.BotOwner) || (u == config.ActiveTeams[t].TeamOwner) {
		privs["SetUserPriv"] = true
		privs["AddUsers"] = true
		privs["KickUsers"] = true
		privs["CreateChannels"] = true
		privs["DeleteChannels"] = true
		privs["SetTopic"] = true
	} else if priv, ok := config.ActiveTeams[t].UserPrivileges[u]; ok {
		privs["SetUserPriv"] = priv.SetUserPriv
		privs["AddUsers"] = priv.AddUsers
		privs["KickUsers"] = priv.KickUsers
		privs["CreateChannels"] = priv.CreateChannels
		privs["DeleteChannels"] = priv.DeleteChannels
		privs["SetTopic"] = priv.SetTopic
	} else {
		privs["SetUserPriv"] = false
		privs["AddUsers"] = false
		privs["KickUsers"] = false
		privs["CreateChannels"] = false
		privs["DeleteChannels"] = false
		privs["SetTopic"] = false
	}

	if len(args) < 3 {
		return parser.CmdOut{}, &cmdError{args[0], "Missing arguments."}
	}

	switch strings.ToLower(args[1]) {
	case "add":
		if !privs["AddUsers"] {
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You do not have permission to add members to *%s*.", message.Msg.Sender.Username, t)}
		}

		team := api.Team{Name: t}
		members := make(map[string]string)
		for _, m := range args[2:] {
			m = strings.TrimPrefix(m, "@")
			m = strings.TrimSuffix(m, ",")
			members[m] = "reader"
		}
		teamAdd := team.AddMembers(members)
		if teamAdd.Error.Message != "" {
			response = teamAdd.Error.Message
		} else {
			response = fmt.Sprintf("%s successfully added to team.", strings.Join(args[2:], ", "))
		}
	case "kick":
		if !privs["KickUsers"] {
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You do not have permission to kick members from *%s*.", message.Msg.Sender.Username, t)}
		}

		team := api.Team{Name: t}
		member := args[2]
		member = strings.TrimPrefix(member, "@")

		if member == config.ActiveTeams[t].TeamOwner {
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You cannot kick the team owner from this team.", message.Msg.Sender.Username)}
		} else if member == config.BotOwner {
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You cannot kick the botOwner from this team.", message.Msg.Sender.Username)}
		} else if member == config.BotUser {
			return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("@%s, You cannot kick the botUser from this team.", message.Msg.Sender.Username)}
		}
		teamKick := team.RemoveMember(member)
		if teamKick.Error.Message != "" {
			response = teamKick.Error.Message
		} else {
			response = fmt.Sprintf("%s successfully kicked from team.", args[2])
		}
	default:
		return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("%s - Invalid argument.", args[1])}
	}

	return parser.CmdOut{response, channel}, nil
}

func init() {
	parser.RegisterCommand("user", "User privilege settings.", true, true, cmdUser)
}
