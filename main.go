package main

import (
	"fmt"
	"os"

	"keybot/api"
	"keybot/config"
)

const(
	ConfigFile = "config.json"
)

func main() {
	c := config.ConfigJSON{}

	// Create default config if none exists
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		c.Write(ConfigFile)
	}

	// Read config file
	c.Read(ConfigFile)

	u := api.Channel{Name: "dxb"}
	message := fmt.Sprintf("Bot owner: %s", c.BotOwner)
	u.SendMessage(message)

	c.BotOwner = "SomeOtherGuy"
	c.Write(ConfigFile)
	message = fmt.Sprintf("New bot owner: %s", c.BotOwner)
	u.SendMessage(message)
}
