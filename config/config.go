package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type KeybotConfig interface {
	Read(filename string)
	Write(filename string)
}

type ConfigJSON struct {
	BotOwner    string                      `json:"botOwner"`
	ActiveTeams map[string]ConfigActiveTeam `json:"activeTeams,omitempty"`
}

type ConfigActiveTeam struct {
	TeamName       string                `json:"teamName"`
	TeamOwner      string                `json:"teamOwner"`
	UserPrivileges []ConfigUserPrivilege `json:"userPrivileges"`
	ActiveChannels []string              `json:"activeChannels"`
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
