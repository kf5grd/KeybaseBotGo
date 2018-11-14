package main

import (
	"bufio"
	"os/exec"

	"keybot/api"
	"keybot/config"
)

const(
	CommandPrefix = "."
)

func main() {
	c := config.ConfigJSON{}

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
		commandHandler(messageIn, &c)
	}
}
