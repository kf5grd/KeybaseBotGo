package api

import (
	"encoding/json"
	"log"
	"os/exec"
)

type TeamAPIResponse struct {
	Result teamAPIResult `json:"result,omitempty"`
	Error  teamAPIError  `json:"error,omitempty"`
}

type teamAPIResult struct {
	ChatSent     bool `json:"chatSent"`
	CreatorAdded bool `json:"creatorAdded"`
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
