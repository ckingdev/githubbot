package githubbot

import (
	"encoding/json"
	"fmt"
)

type PacketType string

const (
	PingReplyType = "ping-reply"
	PingEventType = "ping-event"

	SendType      = "send"
	SendEventType = "send-event"
	SendReplyType = "send-reply"

	NickType      = "nick"
	NickReplyType = "nick-reply"
	NickEventType = "nick-event"

	JoinEventType = "join-event"

	PartEventType = "part-event"

	AuthType = "auth"

	BounceEventType = "bounce-event"
)

type PacketEvent struct {
	ID    string          `json:"id"`
	Type  PacketType      `json:"type"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

type PingEvent struct {
	Time int64 `json:"time"`
	Next int64 `json:"next"`
}

type PingReply struct {
	UnixTime int64 `json:"time,omitempty"`
}

type SendCommand struct {
	Content string `json:"content"`
	Parent  string `json:"parent"`
}

type Message struct {
	ID              string `json:"id"`
	Parent          string `json:"parent"`
	PreviousEditID  string `json:"previous_edit_id,omitempty"`
	Time            int64  `json:"time"`
	Sender          User   `json:"sender"`
	Content         string `json:"content"`
	EncryptionKeyID string `json:"encryption_key_id,omitempty"`
	Edited          int    `json:"edited,omitempty"`
	Deleted         int    `json:"deleted,omitempty"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ServerID  string `json:"server_id"`
	ServerEra string `json:"server_era"`
}

type SendEvent Message
type SendReply Message

type NickCommand struct {
	Name string `json:"name"`
}

type AuthCommand struct {
	Type     string `json:"type"`
	Passcode string `json:"passcode,omitempty"`
}

func (p *PacketEvent) Payload() (interface{}, error) {
	var payload interface{}
	switch p.Type {
	case PingEventType:
		payload = &PingEvent{}
	case SendType:
		payload = &SendCommand{}
	case PingReplyType:
		payload = &PingReply{}
	case AuthType:
		payload = &AuthCommand{}
	case SendEventType:
		payload = &SendEvent{}
	case SendReplyType:
		payload = &SendReply{}
	default:
		return p.Data, fmt.Errorf("Unexpected packet type: %s", p.Type)
	}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}
