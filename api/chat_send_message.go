package api

import (
	"encoding/json"
	"log"
)

type Channel struct {
	IsTeam  bool
	Name    string
	Channel string
}

func (c Channel) SendMessage(message string) ChatAPIResponse {
	/* Send chat message to a user or team */
	var msgJSON chatAPIOut

	// if sending a message to a team but no channel is passed, send to #general
	if c.IsTeam == true && c.Channel == "" {
		c.Channel = "general"
	}
	
	log.Printf(
		"[ChatAPI.SendMessage] [IsTeam: %v] [Name: %s] [Channel: %s] [Message: %s]\n",
		c.IsTeam,
		c.Name,
		c.Channel,
		message,
	)

	msgJSON.Method                             = "send"
	msgJSON.Params.Options.Channel.Name        = c.Name
	msgJSON.Params.Options.Message.Body        = message

	if c.IsTeam {
		msgJSON.Params.Options.Channel.MembersType = "team"
		msgJSON.Params.Options.Channel.TopicName   = c.Channel
	}

	jsonBytes, _ := json.Marshal(msgJSON)
	return SendChatAPI(string(jsonBytes))
}
