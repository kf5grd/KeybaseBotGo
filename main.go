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
			args := parser.GetArgs(message)
			printArgs := ""
			for _, arg := range args[1:] {
				printArgs += fmt.Sprintf("  %s\n", arg)
			}
			fmt.Printf("----\ncommand: %s\nArguments:\n%s\n----", args[0], printArgs)
		}
	}
}
