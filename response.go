package main

type Response struct {
	Body  interface{}  `json:"body"`
	Error interface{}  `json:"error"`
	Type  ResponseType `json:"type"`
}

type ResponseType int

const (
	ERROR ResponseType = iota
	CONNECTED
	TOPIC_ROOMS
	ME_CHANGED_USERNAME
	OTHER_CHANGED_USERNAME
	ME_JOINED_CHAT
	OTHER_JOINED_CHAT
	ME_LEFT_CHAT
	OTHER_LEFT_CHAT
	ME_MESSAGE_SEND
	OTHER_MESSAGE_SEND
	OLD_MESSAGES
)

func (t ResponseType) String() string {
	return []string{
		"ERROR",
		"CONNECTED",
		"TOPIC_ROOMS",
		"ME_CHANGED_USERNAME",
		"OTHER_CHANGED_USERNAME",
		"ME_JOINED_CHAT",
		"OTHER_JOINED_CHAT",
		"ME_LEFT_CHAT",
		"OTHER_LEFT_CHAT",
		"ME_MESSAGE_SEND",
		"OTHER_MESSAGE_SEND",
		"OLD_MESSAGES",
	}[t]
}
