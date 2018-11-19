package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"

	"keybot/api"
	"keybot/commands"
	"keybot/config"
)

type keybaseStatusJSON struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
}

func getKeybaseStatus() keybaseStatusJSON {
	cmd := exec.Command("keybase", "status", "-j")

	cmdOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var retVal keybaseStatusJSON
	json.Unmarshal(cmdOut, &retVal)

	return retVal
}

func main() {
	log.Println("Getting Keybase status...")
	keybaseStatus := getKeybaseStatus()
	if !keybaseStatus.LoggedIn {
		panic("Not logged in to Keybase.")
	}
	
	log.Println("Reading config...")
	c := config.ConfigJSON{}

	// Read config file
	c.Read()
	if c.BotUser != keybaseStatus.Username {
		c.BotUser = keybaseStatus.Username
		c.Write()
	}
	if c.CommandPrefix == "" {
		c.CommandPrefix = "."
		c.Write()
	}

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
		commands.CommandHandler(messageIn, &c)
	}
}
