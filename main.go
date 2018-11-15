package main

import (
	"bufio"
	"log"
	"os/exec"

	"keybot/api"
	"keybot/config"
)

const(
	CommandPrefix = "."
)

func main() {
	log.Println("Reading config...")
	c := config.ConfigJSON{}

	// Read config file
	c.Read()

/*
	c.ActiveTeams = make(map[string]config.ConfigActiveTeam)
	c.ActiveTeams["crbot.public"] = config.ConfigActiveTeam{
		TeamName: "crbot.public",
		TeamOwner: "dxb",
		ActiveChannels: map[string]struct{}{"bots": {}, "test": {}},
	}
	c.ActiveTeams["pho_enix"] = config.ConfigActiveTeam{
		TeamName: "pho_enix",
		TeamOwner: "dxb",
		ActiveChannels: map[string]struct{}{"general": {}},
	}
*/
	c.Write()

	// spawn keybase chat listener and process messages as they come in
	log.Println("Starting chat listener...")
	keybaseListen := exec.Command("keybase", "chat", "api-listen")
	keybaseOutput, _ := keybaseListen.StdoutPipe()
	keybaseListen.Start()
	scanner := bufio.NewScanner(keybaseOutput)
	for scanner.Scan() {
		messageIn := api.ReceiveMessage(scanner.Text())
		commandHandler(messageIn, &c)
	}
}
