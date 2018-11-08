package api

import (
	"encoding/json"
	"log"
)

type ChatAPI interface {
	SendMessage(message string) ChatAPIResponse
}

type Channel struct {
	IsTeam  bool
	Name    string
	Channel string
}

type chatAPIOut struct {
	Method string `json:"method"`
	Params params `json:"params,omitempty"`
}

type params struct {
	Options options `json:"options"`
}

type options struct {
	Channel channel  `json:"channel,omitempty"`
	Message message  `json:"message"`
}

type channel struct {
	Name        string `json:"name"`
	MembersType string `json:"members_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}

type message struct {
	Body string `json:"body"`
}

func (c Channel) SendMessage(message string) ChatAPIResponse {
	var msgJSON chatAPIOut
	log.Printf("[ChatAPI.SendMessage] [IsTeam: %v] [Name: %s] [Channel: %s] [Message: %s]\n", c.IsTeam, c.Name, c.Channel, message)

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
