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

			added := ""
			notadded := ""
			errors := ""
			for _, user := range args[3:] {
				if _, ok := config.Blacklist[user]; ok {
					notadded += fmt.Sprintf("%s, ", user)
				} else {
					switch user {
					case config.BotOwner:
						errors += fmt.Sprintf("  %s: Cannot blacklist bot owner.\n", user)
					case message.Msg.Sender.Username:
						errors += fmt.Sprintf("  %s: Cannot blacklist self.\n", user)
					default:
						added += fmt.Sprintf("%s, ", user)
						config.Blacklist[user] = struct{}{}
					}
				}
			}
			havehas := "has"
			isare := "is"
			response = ""
			if added != "" {
				added = strings.TrimSuffix(added, ", ")
				if strings.Contains(added, ",") {
					havehas = "have"
				}
				response += fmt.Sprintf("%s %s been added to the blacklist.\n", added, havehas)
			}
			if notadded != "" {
				notadded = strings.TrimSuffix(notadded, ", ")
				if strings.Contains(notadded, ",") {
					isare = "are"
				}
				response += fmt.Sprintf("%s %s already on the blacklist.", notadded, isare)
			}
			if errors != "" {
				errors = strings.TrimSuffix(errors, "\n")
				response += "\n*Errors:*\n"
				response += errors
			}
			response = strings.TrimSuffix(response, "\n")
			config.Write()
		case "remove":
			if len(args) < 4 {
				return parser.CmdOut{}, &cmdError{args[0], fmt.Sprintf("`%s` - No users provided.", args[1])}
			}
			if config.Blacklist == nil {
				config.Blacklist = make(map[string]struct{})
				config.Write()
			}

			removed := ""
			notremoved := ""
			for _, user := range args[3:] {
				if _, ok := config.Blacklist[user]; !ok {
					notremoved += fmt.Sprintf("%s, ", user)
				} else {
					removed += fmt.Sprintf("%s, ", user)
					delete(config.Blacklist, user)
				}
			}
			havehas := "has"
			isare := "is"
			response = ""
			if removed != "" {
				removed = strings.TrimSuffix(removed, ", ")
				if strings.Contains(removed, ",") {
					havehas = "have"
				}
				response += fmt.Sprintf("%s %s been removed from the blacklist.\n", removed, havehas)
			}
			if notremoved != "" {
				notremoved = strings.TrimSuffix(notremoved, ", ")
				if strings.Contains(notremoved, ",") {
					isare = "are"
				}
				response += fmt.Sprintf("%s %s not on the blacklist.", notremoved, isare)
			}
			response = strings.TrimSuffix(response, "\n")
			config.Write()
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
	parser.RegisterCommand("user", "User privilege settings.", true, true, cmdUser)
}

func commandHandler(message api.ChatMessageIn, c *config.ConfigJSON) {
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
