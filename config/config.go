package config

import (
	"encoding/json"
)

type KeybotConfig interface {
	Read(filename string) ConfigJSON
	Write(filename string) ConfigJSON
}

type ConfigJSON struct {
	BotOwner    string             `json:"botOwner"`
	ActiveTeams []configActiveTeam `json:"activeTeams"`
}

type configActiveTeam struct {
	TeamOwner      string                `json:"teamOwner"`
	UserPrivileges []configUserPrivilege `json:"userPrivileges"`
	ActiveChannels []string              `json:"activeChannels"`
}

type configUserPrivilege struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (c KeybotConfig) Read(filename string) ConfigJSON {
}

func (c KeybotConfig) Write(filename string) ConfigJSON {
}
