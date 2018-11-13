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
				sendMessage := fmt.Sprintf("Command: %s", args[0])
				if len(args) > 1 {
					sendMessage += ", Args: ["
					for i, arg := range args[1:] {
						sendMessage += fmt.Sprintf("\"%s\"", arg)
						if i != (len(args[1:]) - 1) {
							sendMessage += ", "
						}
					}
					sendMessage += "]"
				}
				m := api.Channel{Name: messageIn.Msg.Channel.Name, Channel: messageIn.Msg.Channel.TopicName, IsTeam: true}
				m.SendMessage(sendMessage)
			}
		}
	}
}
