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

	c.ActiveTeams = make(map[string]config.ConfigActiveTeam)
	c.ActiveTeams["crbot.public"] = config.ConfigActiveTeam{
		TeamName: "crbot.public",
		TeamOwner: "dxb",
		ActiveChannels: []string{"bots", "test"},
	}
	c.Write(ConfigFile)

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
				team := messageIn.Msg.Channel.Name
				channel := messageIn.Msg.Channel.TopicName
				
				isActiveChannel := false
				_, isActiveTeam := c.ActiveTeams[messageIn.Msg.Channel.Name]
				if isActiveTeam {
					for _, ch := range c.ActiveTeams[team].ActiveChannels {
						if ch == channel {
							isActiveChannel = true
						}
					}
				}

				cmd, isCommand := parser.Commands[args[0]]
				if isCommand && isActiveChannel {
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
