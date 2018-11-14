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
	filename := c.Filename
	if filename == "" {
		filename = defaultFilename
	}

	configFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonBytes, _ := ioutil.ReadAll(configFile)

	json.Unmarshal([]byte(jsonBytes), &c)
}

func (c ConfigJSON) Write() {
	filename := c.Filename
	if filename == "" {
		filename = defaultFilename
	}

	configFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
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
