package main

import "github.com/gofiber/websocket/v2"

type Client struct {
	ID         string          `json:"id"`
	Connection *websocket.Conn `json:"-"`
	Username   string          `json:"username"`
}

type Clients struct {
	clients map[string]Client
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[string]Client),
	}
}

func (c *Clients) Add(client Client) {
	c.clients[client.ID] = client
}

func (c *Clients) Remove(clientID string) {
	delete(c.clients, clientID)
}

func (c *Clients) Get(clientID string) (Client, bool) {
	client, ok := c.clients[clientID]
	return client, ok
}
