package main

import (
	"fmt"

	"keybot/api"
)

func main() {
	u := api.Channel{Name: "dxb"}
	t := api.Team{Name: "crbot.public"}

	members := t.ListMembers()

	msg := "```\n"
	for i, member := range members {
		msg += fmt.Sprintf("%d: %s, %s\n", i, member.Username, member.Role)
	}
	msg += "```"
	u.SendMessage(msg)
}
