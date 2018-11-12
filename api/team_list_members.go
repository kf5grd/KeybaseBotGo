package api

import (
	"encoding/json"
	"log"
)

func (t Team) ListMembers() map[string]string {
	var msgJSON teamAPIOut

	log.Printf(
		"[TeamAPI.ListMembers] [Team: %s]",
		t.Name,
	)

	msgJSON.Method              = "list-team-memberships"
	msgJSON.Params.Options.Team = t.Name

	jsonBytes, _ := json.Marshal(msgJSON)
	apiResp := SendTeamAPI(string(jsonBytes))
	retVal := make(map[string]string)
	type members []TeamMember
	
	roles := map[string]members{
		"owner": apiResp.Result.Members.Owners,
		"admin": apiResp.Result.Members.Admins,
		"writer": apiResp.Result.Members.Writers,
		"reader": apiResp.Result.Members.Readers,
	}

	for role, memberList := range roles {
		for _, member := range memberList {
			retVal[member.Username] = role
		}
	}

	return retVal
}
