package main

import (
	"fmt"

	"keybot/api"
)

func main() {
	u := api.Channel{Name: "dxb"}
	t := api.Team{Name: "crbot.public"}

	members := map[string]string{
		"cagingroyals": "admin",
	}

	teamAdd := t.AddMembers(members)
	if teamAdd.Error.Message != "" {
		u.SendMessage(teamAdd.Error.Message)
	} else {
		msg := fmt.Sprintf("Users successfully added to team `%s`.", t.Name)
		u.SendMessage(msg)
	}
}
