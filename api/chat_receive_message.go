package api

import (
	"encoding/json"
	"log"
)

type ChatMessageIn struct {
	Source     string          `json:"source"`
	Msg        chatRecvMessage `json:"msg"`
	Pagination chatPagination  `json:"pagination"`
}

type chatRecvMessage struct {
	ID             int            `json:"id"`
	Channel        chatMsgChannel `json:"channel"`
	Sender         chatMsgSender  `json:"sender"`
	SentAt         int            `json:"sent_at"`
	SentAtMs       int            `json:"sent_at_ms"`
	Content        chatMsgContent `json:"content"`
	Prev           string         `json:"prev"`
	Unread         bool           `json:"unread"`
	ChannelMention string         `json:"channel_mention"`
}

type chatMsgChannel struct{
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	MembersType string `json:"members_type"`
	TopicType   string `json:"topic_type"`
	TopicName   string `json:"topic_name"`
}

type chatMsgSender struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type chatMsgContent struct {
	MsgType string          `json:"type"`
	Edit    chatMsgTypeEdit `json:"edit"`
	Text    chatMsgTypeText `json:"text"`
}

type chatMsgTypeEdit struct {
	MessageID int    `json:"messageID"`
	Body      string `json:"body"`
}

type chatMsgTypeText struct {
	Body string `json:"body"`
}

type chatPagination struct {
	Next     string
	Previous string
	num      int
	last     bool
}

func ReceiveMessage(jsonString string) ChatMessageIn {
	log.Println("[ReceiveMessage]", jsonString)
	var jsonData ChatMessageIn
	json.Unmarshal([]byte(jsonString), &jsonData)
	return jsonData
}
