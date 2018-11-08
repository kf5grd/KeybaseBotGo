package api

import (
	"encoding/json"
	"log"
	"os/exec"
)

type TeamAPI interface {
	ListMembers() TeamAPIResponse
}	

// -- JSON Out to API --
type teamAPIOut struct {
	Method string     `json:"method"`
	Params teamParams `json:"params,omitempty"`
}

type teamParams struct {
	Options teamOptions `json:"options,omitempty"`
}

type teamOptions struct {
	Team string `json:"team,omitempty"`
}

// -- JSON Received back from API --
type TeamAPIResponse struct {
	Result teamAPIResult `json:"result,omitempty"`
	Error  teamAPIError  `json:"error,omitempty"`
}

type teamAPIResult struct {
	ChatSent     bool `json:"chatSent,omitempty"`
	CreatorAdded bool `json:"creatorAdded,omitempty"`

	// list-user-memberships
	Members teamMembers `json:"members,omitempty"`
}

type teamMembers struct {
	owners  []teamMember `json:"owners,omitempty"`
	admins  []teamMember `json:"admins,omitempty"`
	readers []teamMember `json:"readers,omitempty"`
	writers []teamMember `json:"writers,omitempty"`
}

type teamMember struct {
	uv       teamMemberUV `json:"uv,omitempty"`
	Username string       `json:"username,omitempty"`
	FullName string       `json:"fullname,omitempty"`
	needspuk bool         `json:"needsPUK,omitempty"`
	status   int          `json:"status,omitempty"`
}

type teamAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendTeamAPI(jsonData string) TeamAPIResponse {
	log.Println("[SendTeamAPI]", "[out]", jsonData)
	cmd := exec.Command("keybase", "team", "api", "-m", jsonData)

	cmdOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var retVal TeamAPIResponse
	json.Unmarshal(cmdOut, &retVal)

	if retVal.Error.Message != "" {
		log.Printf(
			"[SendTeamAPI] [in] [Error] [Code: %d] [Message: %s]\n",
			retVal.Error.Code,
			retVal.Error.Message,
		)
	} else {
		log.Printf(
			"[SendTeamAPI] [in] [ChatSent: %v] [CreatorAdded: %v]\n",
			retVal.Result.ChatSent,
			retVal.Result.CreatorAdded,
		)
	}

	return retVal
}
