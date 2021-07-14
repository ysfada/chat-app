package main

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type ConnectionStore interface {
	Store(clientID string, conn *websocket.Conn)
	Load(clientID string) (conn *websocket.Conn, ok bool)
	Delete(clientID string)
}

type InMemoryConnectionStore struct {
	sync.Mutex
	connections map[string]*websocket.Conn
}

var _ ConnectionStore = (*InMemoryConnectionStore)(nil)

func NewInMemoryConnectionStore() *InMemoryConnectionStore {
	return &InMemoryConnectionStore{
		connections: map[string]*websocket.Conn{},
	}
}

func (s *InMemoryConnectionStore) Store(clientID string, conn *websocket.Conn) {
	s.Lock()
	s.connections[clientID] = conn
	s.Unlock()
}

func (s *InMemoryConnectionStore) Load(clientID string) (conn *websocket.Conn, ok bool) {
	s.Lock()
	conn, ok = s.connections[clientID]
	s.Unlock()
	return conn, ok
}

func (s *InMemoryConnectionStore) Delete(clientID string) {
	s.Lock()
	delete(s.connections, clientID)
	s.Unlock()
}
