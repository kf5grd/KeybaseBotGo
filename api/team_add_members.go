package api

import (
	"encoding/json"
	"log"
)

func (t Team) AddMembers(members map[string]string) TeamAPIResponse {
	var msgJSON teamAPIOut

	log.Printf(
		"[TeamAPI.AddMembers] [Team: %s] [Members: %v]",
		t.Name,
		members,
	)

	msgJSON.Method = "add-members"
	msgJSON.Params.Options.Team = t.Name

	var usernames []TeamUsername
	for user, role := range members {
		var username TeamUsername
		username.Username = user
		username.Role = role

		usernames = append(usernames, username)
	}
	msgJSON.Params.Options.Usernames = usernames

	jsonBytes, _ := json.Marshal(msgJSON)
	return SendTeamAPI(string(jsonBytes))
}
