package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

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

	// spawn keybase chat listener and process messages as they come in
	keybaseListen := exec.Command("keybase", "chat", "api-listen", "--local")
	keybaseOutput, _ := keybaseListen.StdoutPipe()
	keybaseListen.Start()
	scanner := bufio.NewScanner(keybaseOutput)
	for scanner.Scan() {
		messageIn := api.ReceiveMessage(scanner.Text())
		user := messageIn.Msg.Sender.Username
		message := messageIn.Msg.Content.Text.Body
		fmt.Println(user + ":", message)
	}
}
