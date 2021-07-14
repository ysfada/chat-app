package main

type Request struct {
	ID       string                 `json:"id"`
	ClientID string                 `json:"clientId"`
	Body     map[string]interface{} `json:"body"`
	Type     RequestType            `json:"type"`
}

type RequestType int

const (
	GET_ROOMS RequestType = iota
	CHANGE_USERNAME
	JOIN_CHAT
	LEFT_CHAT
	SEND_MESSAGE
	GET_OLD_MESSAGES
)

func (t RequestType) String() string {
	return []string{
		"GET_ROOMS",
		"CHANGE_USERNAME",
		"JOIN_CHAT",
		"LEFT_CHAT",
		"SEND_MESSAGE",
		"GET_OLD_MESSAGES",
	}[t]
}
