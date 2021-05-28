package main

type Message struct {
	ID       string `json:"id"`
	ClientID string `json:"clientId"`
	RoomID   string `json:"roomId"`
	Message  string `json:"message"`
	Type     string `json:"type"`
}

type Messages struct {
	messages map[string][]Message
}

func NewMessages() *Messages {
	return &Messages{
		messages: make(map[string][]Message),
	}
}

func (m *Messages) MessageCount(roomID string) int {
	return len(m.messages[roomID])
}

func (m *Messages) GetMessages(roomID string) []Message {
	return m.messages[roomID]
}

func (m *Messages) SetMessages(roomID string, messages []Message) {
	m.messages[roomID] = messages
}

func (m *Messages) AppendMessage(roomID string, message Message) {
	m.messages[roomID] = append(m.messages[roomID], message)
}

func (m *Messages) RemoveMessages(roomID string) {
	delete(m.messages, roomID)
}
