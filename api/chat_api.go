package api

import (
	"encoding/json"
	"log"
	"os/exec"
)

type ChatAPIResponse struct {
	Result ChatAPIResult `json:"result"`
}

type ChatAPIResult struct {
	Message    string                `json:"message"`
	RateLimits []ChatAPIResultLimits `json:"ratelimits"`
}

type ChatAPIResultLimits struct {
	Tank     string `json:"tank"`
	Capacity int    `json:"capacity"`
	Reset    int    `json:"reset"`
	Gas      int    `json:"gas"`
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

	log.Printf("[SendChatAPI] [in] [Tank: %s] [Capacity: %d] [Reset: %d] [Gas: %d]", retVal.Result.RateLimits.Tank, retVal.Result.RateLimits.Capacity, retVal.Result.RateLimits.Reset, retVal.Result.RateLimits.Gas)
	return retVal
}
