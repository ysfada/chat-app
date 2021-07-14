package main

import (
	"sync"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type UserStore interface {
	Store(userID string, user User)
	Load(userID string) (user User, ok bool)
	Delete(userID string)
}

type InMemoryUserStore struct {
	sync.Mutex
	users map[string]User
}

var _ UserStore = (*InMemoryUserStore)(nil)

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: map[string]User{},
	}
}

func (s *InMemoryUserStore) Store(userID string, user User) {
	s.Lock()
	s.users[userID] = user
	s.Unlock()
}

func (s *InMemoryUserStore) Load(userID string) (user User, ok bool) {
	s.Lock()
	user, ok = s.users[userID]
	s.Unlock()
	return user, ok
}

func (s *InMemoryUserStore) Delete(userID string) {
	s.Lock()
	delete(s.users, userID)
	s.Unlock()
}
