package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type KeybotConfig interface {
	Read() ConfigJSON
	Write()
}

type ConfigFile struct {
	Path string
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

func (c ConfigFile) Read() ConfigJSON {
	configFile, err := os.Open(c.Path)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonBytes, _ := ioutil.ReadAll(configFile)

	var retVal ConfigJSON
	json.Unmarshal([]byte(jsonBytes), &retVal)

	return retVal
}

func (c ConfigFile) Write() {
}
