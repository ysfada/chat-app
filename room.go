package main

import (
	"github.com/gofiber/websocket/v2"
)

type Rooms struct {
	rooms map[string]map[*websocket.Conn]string
}

func NewRooms() *Rooms {
	return &Rooms{
		rooms: make(map[string]map[*websocket.Conn]string),
	}
}

func (r *Rooms) DeleteRoom(roomID string) {
	delete(r.rooms, roomID)
}

func (r *Rooms) AddConnection(roomID string, clientID string, conn *websocket.Conn) {
	if !r.roomExists(roomID) {
		r.createRoom(roomID)
	}
	r.rooms[roomID][conn] = clientID
}

func (r *Rooms) RemoveConnection(roomID string, conn *websocket.Conn) {
	delete(r.rooms[roomID], conn)
}

func (r *Rooms) Connections(roomID string) map[*websocket.Conn]string {
	return r.rooms[roomID]
}

func (r *Rooms) ConnectionsCount(roomID string) int {
	return len(r.rooms[roomID])
}

func (r *Rooms) roomExists(roomID string) bool {
	if _, ok := r.rooms[roomID]; ok {
		return true
	}
	return false
}

func (r *Rooms) createRoom(roomID string) {
	if !r.roomExists(roomID) {
		r.rooms[roomID] = make(map[*websocket.Conn]string)
	}
}
