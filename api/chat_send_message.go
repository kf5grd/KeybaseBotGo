package api

import (
	"encoding/json"
	"log"
)

type ChatAPI interface {
	SendTeamMsg(team channel, message string) ChatAPIResponse
	SendUserMsg(user, message string) ChatAPIResponse
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

func (j ChatAPIOut) SendTeamMsg(team, channel, message string) ChatAPIResponse {
	log.Printf("[ChatAPI.SendTeamMessage] [Team: %s] [Channel: %s] [Message: %s]\n", team, channel, message)

	j.Method                             = "send"
	j.Params.Options.Channel.MembersType = "team"
	j.Params.Options.Channel.Name        = team
	j.Params.Options.Channel.TopicName   = channel
	j.Params.Options.Message.Body        = message

	jsonBytes, _ := json.Marshal(j)

	return SendChatAPI(string(jsonBytes))
}

func (j ChatAPIOut) SendUserMsg(user, message string) ChatAPIResponse {
	log.Printf("[ChatAPI.SendTeamMessage] [User: %s] [Message: %s]\n", user, message)

	j.Method                      = "send"
	j.Params.Options.Channel.Name = user
	j.Params.Options.Message.Body = message

	jsonBytes, _ := json.Marshal(j)

	return SendChatAPI(string(jsonBytes))
}
