package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	User      *User  `json:"user"`
	RoomID    string `json:"roomId"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type MessageStore interface {
	Count(roomID string) int
	Get(roomID string) []Message
	GetLastN(roomID string, n int, firstMsgID ...string) []Message
	Set(roomID string, messages []Message)
	Append(roomID string, message Message)
}

type InMemoryMessageStore struct {
	sync.Mutex
	messages map[string][]Message
}

var _ MessageStore = (*InMemoryMessageStore)(nil)

func NewInMemoryMessageStore() *InMemoryMessageStore {
	m := &InMemoryMessageStore{
		messages: map[string][]Message{},
	}

	for i := 1; i < 201; i++ {
		id, err := uuid.NewRandom()
		if err != nil {
			continue
		}
		m.messages["77dac06c-bb59-4854-8b4b-928d078454cc"] = append(m.messages["77dac06c-bb59-4854-8b4b-928d078454cc"],
			Message{
				ID:        id.String(),
				UserID:    fmt.Sprintf("initial %d", i),
				RoomID:    "77dac06c-bb59-4854-8b4b-928d078454cc",
				Message:   fmt.Sprintf("initial %d", i),
				Timestamp: time.Now().Unix() * 1000,
			})
		m.messages["fe09b952-7690-4978-96cf-5a5c8e74ecaf"] = append(m.messages["fe09b952-7690-4978-96cf-5a5c8e74ecaf"],
			Message{
				ID:        id.String(),
				UserID:    fmt.Sprintf("initial %d", i),
				RoomID:    "fe09b952-7690-4978-96cf-5a5c8e74ecaf",
				Message:   fmt.Sprintf("initial %d", i),
				Timestamp: time.Now().Unix() * 1000,
			})
	}
	return m
}

func (m *InMemoryMessageStore) Count(roomID string) int {
	m.Lock()
	l := len(m.messages[roomID])
	m.Unlock()
	return l
}

func (m *InMemoryMessageStore) Get(roomID string) []Message {
	m.Lock()
	messages := m.messages[roomID]
	m.Unlock()
	return messages
}

func (m *InMemoryMessageStore) GetLastN(roomID string, n int, firstMsgID ...string) []Message {
	var messages []Message
	if len(firstMsgID) > 0 {
		i := m.indexOf(roomID, firstMsgID[0])
		if i < 0 {
			s := len(m.messages[roomID]) - n
			if s < 0 {
				s = 0
			}
			m.Lock()
			messages = m.messages[roomID][s:]
			m.Unlock()
		} else {
			s := i - n
			if s < 0 {
				s = 0
			}
			m.Lock()
			messages = m.messages[roomID][s:i]
			m.Unlock()
		}
	} else {
		s := len(m.messages[roomID]) - n
		if s < 0 {
			s = 0
		}
		m.Lock()
		messages = m.messages[roomID][s:]
		m.Unlock()
	}

	return messages
}

func (m *InMemoryMessageStore) Set(roomID string, messages []Message) {
	m.Lock()
	m.messages[roomID] = messages
	m.Unlock()
}

func (m *InMemoryMessageStore) Append(roomID string, message Message) {
	m.Lock()
	m.messages[roomID] = append(m.messages[roomID], message)
	m.Unlock()
}

func (m *InMemoryMessageStore) indexOf(roomID string, msgID string) int {
	m.Lock()
	for i, msg := range m.messages[roomID] {
		if msg.ID == msgID {
			m.Unlock()
			return i
		}
	}
	m.Unlock()
	return -1 //not found.
}
