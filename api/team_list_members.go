package api

import (
	"encoding/json"
	"log"
)

type Member struct {
	Username string
	FullName string
	Role     string
}

func (t Team) ListMembers() []Member {
	var msgJSON teamAPIOut

	log.Printf(
		"[TeamAPI.ListMembers] [Team: %s]",
		t.Name,
	)

	msgJSON.Method              = "list-team-memberships"
	msgJSON.Params.Options.Team = t.Name

	jsonBytes, _ := json.Marshal(msgJSON)
	apiResp := SendTeamAPI(string(jsonBytes))
	var retVal []Member
	type members []TeamMember
	
	roles := map[string]members{
		"owner": apiResp.Result.Members.Owners,
		"admin": apiResp.Result.Members.Admins,
		"writer": apiResp.Result.Members.Writers,
		"reader": apiResp.Result.Members.Readers,
	}

	for role, memberList := range roles {
		for _, member := range memberList {
			var newMember = Member{
				Username: member.Username,
				FullName: member.FullName,
				Role:     role,
			}
			retVal = append(retVal, newMember)
		}
	}

	return retVal
}
