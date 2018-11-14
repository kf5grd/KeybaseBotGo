package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const defaultFilename string = "config.json"

type KeybotConfig interface {
	Read()
	Write()
}

type ConfigJSON struct {
	Filename    string                      `json:"-"`
	BotOwner    string                      `json:"botOwner"`
	ActiveTeams map[string]ConfigActiveTeam `json:"activeTeams,omitempty"`
	Blacklist   map[string]struct{}         `json:"blacklist,omitempty"`
}

type ConfigActiveTeam struct {
	TeamName       string                `json:"teamName"`
	TeamOwner      string                `json:"teamOwner"`
	UserPrivileges []ConfigUserPrivilege `json:"userPrivileges"`
	ActiveChannels map[string]struct{}   `json:"activeChannels"`
}

type ConfigUserPrivilege struct {
	Username       string `json:"username"`
	SetUserPriv    bool   `json:"setUserPriv"`
	AddUsers       bool   `json:"addUsers"`
	KickUsers      bool   `json:"kickUsers"`
	CreateChannels bool   `json:"createChannels"`
	DeleteChannels bool   `json:"deleteChannels"`
	SetTopic       bool   `json:"setTopic"`
}

func (c *ConfigJSON) Read() {
	if c.Filename == "" {
		c.Filename = defaultFilename
	}

	// Create default config if none exists
	if _, err := os.Stat(c.Filename); os.IsNotExist(err) {
		c.Write()
	}

	configFile, err := os.Open(c.Filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonBytes, _ := ioutil.ReadAll(configFile)

	json.Unmarshal([]byte(jsonBytes), &c)
}

func (c ConfigJSON) Write() {
	if c.Filename == "" {
		c.Filename = defaultFilename
	}

	configFile, err := os.OpenFile(c.Filename, os.O_WRONLY|os.O_CREATE, 0644)
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
