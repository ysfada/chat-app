package main

import (
	"sync"
)

type RoomType int

const (
	TopicRoom RoomType = iota
	UserRoom
)

type Room struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Type  RoomType `json:"type"`
	Users []string `json:"-"`
}

type RoomStore interface {
	Create(roomID string, roomName string, roomType RoomType)
	Join(roomID string, userID string) bool
	Leave(roomID string, userID string)
	Users(roomID string) []string
	Room(roomID string) (room Room, ok bool)
	Rooms(includeUserRoom ...bool) []Room
	UserJoinedTo(userID string) (room Room, ok bool)
}

type InMemoryRoomStore struct {
	sync.Mutex
	rooms map[string]Room
}

var _ RoomStore = (*InMemoryRoomStore)(nil)

func NewInMemoryRoomStore() *InMemoryRoomStore {
	r := &InMemoryRoomStore{
		rooms: map[string]Room{
			"09e9a18a-519f-45d8-80fa-238ef384e4b4": {
				ID:    "09e9a18a-519f-45d8-80fa-238ef384e4b4",
				Name:  "Football",
				Type:  TopicRoom,
				Users: []string{},
			},
			"77dac06c-bb59-4854-8b4b-928d078454cc": {
				ID:    "77dac06c-bb59-4854-8b4b-928d078454cc",
				Name:  "Animals",
				Type:  TopicRoom,
				Users: []string{},
			},
			"405608b0-e2cf-4d30-a106-66365f69f8cb": {
				ID:    "405608b0-e2cf-4d30-a106-66365f69f8cb",
				Name:  "Sports",
				Type:  TopicRoom,
				Users: []string{},
			},
			"9c61b8a5-9bef-4232-9979-946386acefe1": {
				ID:    "9c61b8a5-9bef-4232-9979-946386acefe1",
				Name:  "Politics",
				Type:  TopicRoom,
				Users: []string{},
			},
			"f268891a-9b16-4003-9bc5-82299e41ff8c": {
				ID:    "f268891a-9b16-4003-9bc5-82299e41ff8c",
				Name:  "Social",
				Type:  TopicRoom,
				Users: []string{},
			},
			"c854ac7f-b8d9-4958-b8ae-9b16e105717e": {
				ID:    "c854ac7f-b8d9-4958-b8ae-9b16e105717e",
				Name:  "Cryptocurrency",
				Type:  TopicRoom,
				Users: []string{},
			},
			"11424147-cb85-42ff-8da4-63558d134173": {
				ID:    "11424147-cb85-42ff-8da4-63558d134173",
				Name:  "Relationship",
				Type:  TopicRoom,
				Users: []string{},
			},
			"e7019c42-ab25-4f54-943a-2cd1e471d085": {
				ID:    "e7019c42-ab25-4f54-943a-2cd1e471d085",
				Name:  "Programming",
				Type:  TopicRoom,
				Users: []string{},
			},
			"d5da33df-c9fe-488d-bafd-2b6579c07a9e": {
				ID:    "d5da33df-c9fe-488d-bafd-2b6579c07a9e",
				Name:  "Education",
				Type:  TopicRoom,
				Users: []string{},
			},
			"37f63c5d-7655-46b2-b0e9-8fbdef0f4795": {
				ID:    "37f63c5d-7655-46b2-b0e9-8fbdef0f4795",
				Name:  "Marketing",
				Type:  TopicRoom,
				Users: []string{},
			},
			"6b902675-803c-4d42-b8f4-b9c8cce49282": {
				ID:    "6b902675-803c-4d42-b8f4-b9c8cce49282",
				Name:  "Theaters",
				Type:  TopicRoom,
				Users: []string{},
			},
			"fe09b952-7690-4978-96cf-5a5c8e74ecaf": {
				ID:    "fe09b952-7690-4978-96cf-5a5c8e74ecaf",
				Name:  "Books",
				Type:  TopicRoom,
				Users: []string{},
			},
			"f09c1052-7604-40b4-b8dc-8fe239d94dd4": {
				ID:    "f09c1052-7604-40b4-b8dc-8fe239d94dd4",
				Name:  "TV",
				Type:  TopicRoom,
				Users: []string{},
			},
			"71ff157b-1336-4778-9bca-50d54e0ea3b7": {
				ID:    "71ff157b-1336-4778-9bca-50d54e0ea3b7",
				Name:  "Movies",
				Type:  TopicRoom,
				Users: []string{},
			},
			"8a8276f0-6921-43df-a248-a316d8523a66": {
				ID:    "8a8276f0-6921-43df-a248-a316d8523a66",
				Name:  "Random",
				Type:  TopicRoom,
				Users: []string{},
			},
		},
	}
	return r
}

func (r *InMemoryRoomStore) Create(roomID string, roomName string, roomType RoomType) {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if !ok {
		r.Lock()
		r.rooms[roomID] = Room{
			ID:    roomID,
			Name:  roomName,
			Type:  roomType,
			Users: []string{},
		}
		r.Unlock()
	}
}

func (r *InMemoryRoomStore) Join(roomID string, userID string) bool {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if ok {
		r.Lock()
		tmp := r.rooms[roomID]
		tmp.Users = append(tmp.Users, userID)
		r.rooms[roomID] = tmp
		r.Unlock()
	}
	return ok
}

func (r *InMemoryRoomStore) Leave(roomID string, userID string) {
	if roomID == "" { // TODO: find a better way to handle unknown roomId situation
		r.Lock()
		for id, room := range r.rooms {
			for i, cid := range room.Users {
				if userID == cid {
					tmp := r.rooms[id]
					tmp.Users = append(tmp.Users[:i], tmp.Users[i+1:]...)
					r.rooms[id] = tmp
					break // TODO: this break assumes users can join one room at a time
				}
			}
		}
		r.Unlock()
	} else {
		r.Lock()
		for i, cid := range r.rooms[roomID].Users {
			if userID == cid {
				tmp := r.rooms[roomID]
				tmp.Users = append(tmp.Users[:i], tmp.Users[i+1:]...)
				r.rooms[roomID] = tmp
				break // TODO: this break assumes users can join one room at a time
			}
		}
		r.Unlock()
	}
}

func (r *InMemoryRoomStore) Users(roomID string) []string {
	r.Lock()
	room := r.rooms[roomID]
	r.Unlock()

	return room.Users
}

func (r *InMemoryRoomStore) Room(roomID string) (room Room, ok bool) {
	r.Lock()
	room, ok = r.rooms[roomID]
	r.Unlock()
	return room, ok
}

func (r *InMemoryRoomStore) Rooms(includeUserRoom ...bool) []Room {
	incAll := false
	if len := len(includeUserRoom); len > 0 && includeUserRoom[0] {
		incAll = true
	}
	var rooms []Room
	r.Lock()
	for _, room := range r.rooms {
		if room.Type == UserRoom && !incAll {
			continue
		}

		rooms = append(rooms, room)
	}
	r.Unlock()
	return rooms
}

func (r *InMemoryRoomStore) UserJoinedTo(userID string) (room Room, ok bool) {
	var rm Room
	r.Lock()
	for _, room := range r.rooms {
		for _, cid := range room.Users {
			if cid == userID {
				rm = room
				r.Unlock()
				return rm, true
			}
		}
	}
	r.Unlock()
	return rm, false
}
