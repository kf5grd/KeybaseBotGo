package api

import (
	"encoding/json"
	"log"
	"os/exec"
)

type TeamAPI interface {
	AddMembers(members map[string]string) TeamAPIResponse
	ListMembers() TeamAPIResponse
}	

type Team struct {
	Name string
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
	Team      string         `json:"team,omitempty"`
	Emails    []TeamEmail    `json:"emails,omitempty"`
	Usernames []TeamUsername `json:"usernames,omitempty"`
}

type TeamEmail struct {
	Email string `json:"email"`
	Role  string `json:"role,omitempty"`
}

type TeamUsername struct {
	Username string `json:"username"`
	Role     string `json:"role,omitempty"`
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
	Owners  []TeamMember `json:"owners,omitempty"`
	Admins  []TeamMember `json:"admins,omitempty"`
	Readers []TeamMember `json:"readers,omitempty"`
	Writers []TeamMember `json:"writers,omitempty"`
}

type TeamMember struct {
	uv       teamMemberUV `json:"uv,omitempty"`
	Username string       `json:"username,omitempty"`
	FullName string       `json:"fullname,omitempty"`
	needspuk bool         `json:"needsPUK,omitempty"`
	status   int          `json:"status,omitempty"`
}

type teamMemberUV struct {
	uid         string `json:"uid,omitempty"`
	eldestSeqno int    `json:"eldestSeqno,omitempty"`
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
