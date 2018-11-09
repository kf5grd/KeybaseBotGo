package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"fmt"
)

type KeybotConfig interface {
	Read(filename string)
	Write(filename string)
}

type ConfigJSON struct {
	BotOwner    string             `json:"botOwner"`
	ActiveTeams []configActiveTeam `json:"activeTeams"`
}

type configActiveTeam struct {
	TeamName       string                `json:"teamName"`
	TeamOwner      string                `json:"teamOwner"`
	UserPrivileges []configUserPrivilege `json:"userPrivileges"`
	ActiveChannels []string              `json:"activeChannels"`
}

type configUserPrivilege struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (c *ConfigJSON) Read(filename string) {
	configFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonBytes, _ := ioutil.ReadAll(configFile)

	json.Unmarshal([]byte(jsonBytes), &c)
}

func (c ConfigJSON) Write(filename string) {
	configFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	jsonBytes, _ := json.MarshalIndent(c, "", "  ")
	_, err = configFile.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}
