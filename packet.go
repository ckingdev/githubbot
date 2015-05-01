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
	default:
		return p.Data, fmt.Errorf("Unexpected packet type: %s", p.Type)
	}
	err := json.Unmarshal(p.Data, &payload)
	return payload, err
}
