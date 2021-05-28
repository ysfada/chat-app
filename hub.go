package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type BaseResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Type     string `json:"type"`
}

type MessagesResponse struct {
	Messages []BaseResponse `json:"messages"`
	Type     string         `json:"type"`
}

type HubOptions struct {
	MaxMessage int `json:"maxMessage"`
}

type Hub struct {
	Register   chan *websocket.Conn
	Broadcast  chan *Message
	NewUser    chan *Message
	Messages   chan *Message
	Unregister chan *websocket.Conn
	Options    *HubOptions
	room       *Rooms
	client     *Clients
	message    *Messages
}

func NewHub() *Hub {
	return &Hub{
		Register:   make(chan *websocket.Conn),
		Broadcast:  make(chan *Message),
		NewUser:    make(chan *Message),
		Messages:   make(chan *Message),
		Unregister: make(chan *websocket.Conn),
		Options: &HubOptions{
			MaxMessage: 10,
		},
		room:    NewRooms(),
		client:  NewClients(),
		message: NewMessages(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case connection := <-h.Register:
			h.register(connection)

		case message := <-h.Broadcast:
			h.broadcast(message)

		case message := <-h.NewUser:
			h.newUser(message)

		case message := <-h.Messages:
			h.messages(message)

		case connection := <-h.Unregister:
			h.unregister(connection)
		}
	}
}

func (h *Hub) Upgrade(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client
	// requested upgrade to the WebSocket protocol.
	if websocket.IsWebSocketUpgrade(c) {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		c.Locals("ClientID", uuid.String())
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Hub) Handler(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		h.Unregister <- c
		c.Close()
	}()

	// Register the client
	h.Register <- c

	for {
		var msg Message
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("UnexpectedCloseError: %#v\n", err)
			} else {
				log.Printf("ExpectedCloseError: %#v\n", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		uuid, err := uuid.NewRandom()
		if err != nil {
			return
		}
		msg.ID = uuid.String()
		msg.ClientID = c.Locals("ClientID").(string)
		msg.RoomID = c.Params("roomID")

		switch msg.Type {
		case "MESSAGE":
			h.Broadcast <- &msg
		case "NEW_USER":
			h.NewUser <- &msg
		case "MESSAGES":
			h.Messages <- &msg
		default:
			log.Printf("Unknown message type received: %#v\n", msg)
			return // Calls the deferred function, i.e. closes the connection on error
		}
	}
}

func (h *Hub) register(connection *websocket.Conn) {
	client := Client{
		ID:         connection.Locals("ClientID").(string),
		Username:   "<unset>",
		Connection: connection,
	}
	h.client.Add(client)

	roomID := connection.Params("roomID")
	h.room.AddConnection(roomID, client.ID, connection)
}

func (h *Hub) broadcast(message *Message) {
	if h.message.MessageCount(message.RoomID) >= h.Options.MaxMessage {
		h.message.SetMessages(message.RoomID, h.message.GetMessages(message.RoomID)[1:])
	}

	h.message.AppendMessage(message.RoomID, *message)

	for conn, clientID := range h.room.Connections(message.RoomID) {
		if message.ClientID == conn.Locals("ClientID") {
			continue
		}
		client, ok := h.client.Get(message.ClientID)
		if !ok {
			continue
		}
		if err := conn.WriteJSON(BaseResponse{
			Username: client.Username,
			Message:  message.Message,
			Type:     "MESSAGE",
		}); err != nil {
			log.Printf("(roomID: %s; clientID: %s) write error: %s\n", message.RoomID, clientID, err.Error())
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()
			h.room.RemoveConnection(message.RoomID, conn)
			h.client.Remove(clientID)
		}
	}
}

func (h *Hub) newUser(message *Message) {
	client, ok := h.client.Get(message.ClientID)
	if !ok {
		return // Calls the deferred function, i.e. closes the connection on error
	}
	client.Username = message.Message
	h.client.Add(client)

	for conn, clientID := range h.room.Connections(message.RoomID) {
		if message.ClientID == clientID {
			continue
		}
		if err := conn.WriteJSON(BaseResponse{
			Username: client.Username,
			Message:  "joined",
			Type:     "JOIN",
		}); err != nil {
			log.Printf("(roomID: %s; clientID: %s) write error: %s\n", message.RoomID, clientID, err.Error())
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			conn.Close()
			h.room.RemoveConnection(message.RoomID, conn)
			h.client.Remove(clientID)
		}
	}
}

func (h *Hub) messages(message *Message) {
	client, ok := h.client.Get(message.ClientID)
	if !ok {
		return // Calls the deferred function, i.e. closes the connection on error
	}

	var messages MessagesResponse = MessagesResponse{
		Type: "MESSAGES",
	}
	for _, msg := range h.message.GetMessages(message.RoomID) {
		username := "<removed>"
		if client, ok := h.client.Get(msg.ClientID); ok {
			username = client.Username
		}
		messages.Messages = append(messages.Messages, BaseResponse{
			Username: username,
			Message:  msg.Message,
			Type:     "MESSAGE",
		})
	}

	if err := client.Connection.WriteJSON(messages); err != nil {
		log.Printf("(roomID: %s; clientID: %s) write error: %s\n", message.RoomID, client.ID, err.Error())
		client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
		client.Connection.Close()
		h.room.RemoveConnection(message.RoomID, client.Connection)
		h.client.Remove(client.ID)
	}
}

func (h *Hub) unregister(connection *websocket.Conn) {
	roomID := connection.Params("roomID")
	username := "<removed>"
	if client, ok := h.client.Get(connection.Locals("ClientID").(string)); ok {
		username = client.Username
		h.room.RemoveConnection(roomID, client.Connection)
		h.client.Remove(client.ID)
	}

	if h.room.ConnectionsCount(roomID) > 0 {
		for conn, clientID := range h.room.Connections(roomID) {
			if err := conn.WriteJSON(BaseResponse{
				Username: username,
				Message:  "left",
				Type:     "LEFT",
			}); err != nil {
				log.Printf("(roomID: %s; clientID: %s) write error: %s\n", roomID, clientID, err.Error())
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				conn.Close()
				h.room.RemoveConnection(roomID, conn)
				h.client.Remove(clientID)
			}
		}
	} else {
		h.message.RemoveMessages(roomID)
		h.room.DeleteRoom(roomID)
	}
}
