package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"keybot/api"
	"keybot/config"
	"keybot/parser"
)

const(
	ConfigFile = "config.json"
	CommandPrefix = "."
)

func cmdTest(args []string, message api.ChatMessageIn, config config.ConfigJSON) (parser.CmdOut, error) {
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

func init() {
	parser.RegisterCommand("ping", "Responds with 'pong'", true, true, cmdTest)
}

func main() {
	c := config.ConfigJSON{}

	// Create default config if none exists
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		c.Write(ConfigFile)
	}

	// Read config file
	c.Read(ConfigFile)

	// spawn keybase chat listener and process messages as they come in
	keybaseListen := exec.Command("keybase", "chat", "api-listen", "--local")
	keybaseOutput, _ := keybaseListen.StdoutPipe()
	keybaseListen.Start()
	scanner := bufio.NewScanner(keybaseOutput)
	for scanner.Scan() {
		messageIn := api.ReceiveMessage(scanner.Text())
		message := strings.TrimSpace(messageIn.Msg.Content.Text.Body)
		if strings.HasPrefix(message, CommandPrefix) {
			message = strings.TrimPrefix(message, CommandPrefix)
			args, err := parser.GetArgs(message)
			if err != nil {
				fmt.Println(err)
			} else {
				if cmd, ok := parser.Commands[args[0]]; ok {
					m, err := cmd.CmdFunc(args, messageIn, c)
					if err != nil {
						fmt.Println(err)
					} else {
						m.Channel.SendMessage(m.Message)
					}
				}
			}
		}
	}
}
