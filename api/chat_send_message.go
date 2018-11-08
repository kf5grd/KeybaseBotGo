package api

import (
	"encoding/json"
	"log"
)

type ChatAPI interface {
	SendMessage(message string) ChatAPIResponse
}

type User struct {
	IsTeam  bool
	Name    string
	Channel string
}

type ChatAPIOut struct {
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Params struct {
	Options Options `json:"options"`
}

type Options struct {
	Channel Channel  `json:"channel"`
	Message Message  `json:"message"`
}

type Channel struct {
	Name        string `json:"name"`
	MembersType string `json:"members_type"`
	TopicName   string `json:"topic_name"`
}

type Message struct {
	Body string `json:"body"`
}

func (u User) SendMessage(message string) ChatAPIResponse {
	var msgJSON ChatAPIOut
	log.Printf("[ChatAPI.SendMessage] [IsTeam: %v] [Name: %s] [Channel: %s] [Message: %s]\n", u.IsTeam, u.Name, u.Channel, message)

	msgJSON.Method                             = "send"
	msgJSON.Params.Options.Channel.Name        = u.Name
	msgJSON.Params.Options.Message.Body        = message

	if u.IsTeam {
		msgJSON.Params.Options.Channel.MembersType = "team"
		msgJSON.Params.Options.Channel.TopicName   = u.Channel
	}

	jsonBytes, _ := json.Marshal(msgJSON)

	return SendChatAPI(string(jsonBytes))
}
