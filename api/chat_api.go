package api

import (
	"encoding/json"
	"log"
	"os/exec"
)

type ChatAPIResponse struct {
	Result chatAPIResult `json:"result,omitempty"`
	Error  chatAPIError  `json:"error,omitempty"`
}

type chatAPIResult struct {
	Message    string                `json:"message"`
	RateLimits []chatAPIResultLimits `json:"ratelimits"`
}

type chatAPIResultLimits struct {
	Tank     string `json:"tank"`
	Capacity int    `json:"capacity"`
	Reset    int    `json:"reset"`
	Gas      int    `json:"gas"`
}

type chatAPIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendChatAPI(jsonData string) ChatAPIResponse {
	log.Println("[SendChatAPI]","[out]", jsonData)
	cmd := exec.Command("keybase", "chat", "api", "-m", jsonData)

	cmdOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var retVal ChatAPIResponse
	json.Unmarshal(cmdOut, &retVal)

	if retVal.Error.Message != "" {
		log.Printf(
			"[SendChatAPI] [in] [Error] [Code: %d] [Message: %s]\n",
			retVal.Error.Code,
			retVal.Error.Message,
		)
	} else {
		log.Printf(
			"[SendChatAPI] [in] [Tank: %s] [Capacity: %d] [Reset: %d] [Gas: %d]",
			retVal.Result.RateLimits[0].Tank,
			retVal.Result.RateLimits[0].Capacity,
			retVal.Result.RateLimits[0].Reset,
			retVal.Result.RateLimits[0].Gas,
		)
	}

	return retVal
}
