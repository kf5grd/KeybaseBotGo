package main

import (
	"bufio"
	"os"
	"os/exec"

	"keybot/api"
	"keybot/config"
)

const(
	CommandPrefix = "."
)

func main() {
	c := config.ConfigJSON{}

	// Create default config if none exists
	if _, err := os.Stat(c.Filename); os.IsNotExist(err) {
		c.Write()
	}

	// Read config file
	c.Read()

	c.ActiveTeams = make(map[string]config.ConfigActiveTeam)
	c.ActiveTeams["crbot.public"] = config.ConfigActiveTeam{
		TeamName: "crbot.public",
		TeamOwner: "dxb",
		ActiveChannels: map[string]struct{}{"bots": {}, "test": {}},
	}
	c.Write()

	// spawn keybase chat listener and process messages as they come in
	keybaseListen := exec.Command("keybase", "chat", "api-listen", "--local")
	keybaseOutput, _ := keybaseListen.StdoutPipe()
	keybaseListen.Start()
	scanner := bufio.NewScanner(keybaseOutput)
	for scanner.Scan() {
		messageIn := api.ReceiveMessage(scanner.Text())
		commandHandler(messageIn, c)
	}
}
