package commands

import (
	"fmt"
	"strings"

	"keybot/api"
	"keybot/config"
	"keybot/parser"
)

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
	parser.RegisterCommand("config", "Get and set config values.", true, true, cmdConfig)
}
