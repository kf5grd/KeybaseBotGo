package api

import (
	"encoding/json"
	"log"
)

type Team struct {
	Name string
}

type TeamMember struct {
	Username string
	FullName string
	Role     string
}

func (t Team) ListMembers() []TeamMember {
	var msgJSON teamAPIOut

	log.Printf(
		"[TeamAPI.ListMembers] [Team: %s]",
		t.Name,
	)

	msgJSON.Method              = "list-team-memberships"
	msgJSON.Params.Options.Team = t.Name

	jsonBytes, _ := json.Marshal(msgJSON)
	apiResp := SendTeamAPI(string(jsonBytes))
}
