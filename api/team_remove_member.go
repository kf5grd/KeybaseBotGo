package api

import (
	"encoding/json"
	"log"
)

func (t Team) RemoveMember(member string) TeamAPIResponse {
	var msgJSON teamAPIOut

	log.Printf(
		"[TeamAPI.AddMembers] [Team: %s] [Member: %s]",
		t.Name,
		member,
	)

	msgJSON.Method                   = "remove-member"
	msgJSON.Params.Options.Team      = t.Name
	msgJSON.Params.Options.Username = member

	jsonBytes, _ := json.Marshal(msgJSON)
	return SendTeamAPI(string(jsonBytes))
}
