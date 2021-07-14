package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type HubOptions struct {
	MaxSavedMessage    int
	MaxReturnedMessage int
}

type Hub struct {
	Register       chan *websocket.Conn
	Unregister     chan *websocket.Conn
	GetRooms       chan *websocket.Conn
	ChangeUsername chan *Request
	JoinChat       chan *Request
	LeaveChat      chan *Request
	SendMessage    chan *Request
	OldMessages    chan *Request
	Options        *HubOptions
	connection     ConnectionStore
	user           UserStore
	room           RoomStore
	message        MessageStore
}

func (h *Hub) Defaults() {
	h.Options = &HubOptions{
		MaxSavedMessage:    500,
		MaxReturnedMessage: 20,
	}
	h.connection = NewInMemoryConnectionStore()
	h.user = NewInMemoryUserStore()
	h.room = NewInMemoryRoomStore()
	h.message = NewInMemoryMessageStore()
}

func NewHub() *Hub {
	return &Hub{
		Register:       make(chan *websocket.Conn),
		Unregister:     make(chan *websocket.Conn),
		GetRooms:       make(chan *websocket.Conn),
		ChangeUsername: make(chan *Request),
		JoinChat:       make(chan *Request),
		LeaveChat:      make(chan *Request),
		SendMessage:    make(chan *Request),
		OldMessages:    make(chan *Request),
	}
}

func (h *Hub) Upgrade(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client requested upgrade to the WebSocket protocol.
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

func (h *Hub) Handler(conn *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		h.Unregister <- conn

		if err := conn.Close(); err != nil {
			log.Printf("%#v\n", err)
		}
	}()

	// Register the client
	h.Register <- conn

	for {
		// Read incomming message
		var request Request
		if err := conn.ReadJSON(&request); err != nil {
			if e := h.error(conn, fiber.ErrBadRequest); e != nil {
				return // Calls the deferred function, i.e. closes the connection on error
			}
			continue // Continues on to next request
		}

		// Generate request id
		if id, err := uuid.NewRandom(); err != nil {
			if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
				return // Calls the deferred function, i.e. closes the connection on error
			}
			continue // Continues on to next request
		} else {
			request.ID = id.String()
		}

		// Set ClientID to request
		{
			tmp := conn.Locals("ClientID")
			if str, ok := tmp.(string); ok {
				if id, err := uuid.Parse(str); err == nil {
					request.ClientID = id.String()
				} else {
					if err := h.error(conn, fiber.ErrInternalServerError); err != nil {
						return // Calls the deferred function, i.e. closes the connection on error
					}
					continue // Continues on to next request
				}
			} else {
				if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
					return // Calls the deferred function, i.e. closes the connection on error
				}
				continue // Continues on to next request
			}
		}

		// Handle incomming request base of its type
		switch request.Type {

		case GET_ROOMS:
			h.GetRooms <- conn

		case CHANGE_USERNAME:
			h.ChangeUsername <- &request

		case JOIN_CHAT:
			h.JoinChat <- &request

		case LEFT_CHAT:
			h.LeaveChat <- &request

		case SEND_MESSAGE:
			h.SendMessage <- &request

		case GET_OLD_MESSAGES:
			h.OldMessages <- &request

		default:
			if e := h.error(conn, fiber.ErrBadRequest); e != nil {
				return // Calls the deferred function, i.e. closes the connection on error
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.register(conn)

		case conn := <-h.Unregister:
			h.unregister(conn)

		case conn := <-h.GetRooms:
			h.get_rooms(conn)

		case req := <-h.ChangeUsername:
			h.change_username(req)

		case req := <-h.JoinChat:
			h.join_chat(req)

		case req := <-h.LeaveChat:
			h.leave_chat(req)

		case req := <-h.SendMessage:
			h.send_message(req)

		case req := <-h.OldMessages:
			h.old_messages(req)
		}
	}
}

func (h *Hub) register(conn *websocket.Conn) {
	// Read ClientID
	var clientID string
	{
		tmp := conn.Locals("ClientID")
		if str, ok := tmp.(string); ok && len(str) > 0 {
			clientID = str
		} else {
			h.unregister(conn)
			return
		}
	}

	// Create user
	user := User{
		ID:       clientID,
		Username: clientID,                      // TODO: generate random username
		Avatar:   "https://picsum.photos/56/56", // TODO: create a default avatar to use
		// Avatar:   "https://thispersondoesnotexist.com/image", // TODO: create a default avatar to use
	}
	// Store connection
	h.connection.Store(user.ID, conn)
	// Store user
	h.user.Store(user.ID, user)
	// Create new room for user
	h.room.Create(user.ID, user.Username, UserRoom)

	res := Response{
		Body: map[string]interface{}{
			"message": "connection successful",
			"data":    &user,
		},
		Type: CONNECTED,
	}
	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) unregister(conn *websocket.Conn) {
	// Read ClientID
	var clientID string
	{
		tmp := conn.Locals("ClientID")
		if str, ok := tmp.(string); ok && len(str) > 0 {
			clientID = str
		} else {
			h.error(conn, fiber.ErrInternalServerError)
			return
		}
	}

	// Load user
	user, ok := h.user.Load(clientID)
	if !ok {
		user.ID = "<removed>"
		user.Username = "<removed>"
	}

	// If user joined a chat room than get room id
	var roomID string
	{
		if room, ok := h.room.UserJoinedTo(user.ID); ok {
			roomID = room.ID
		}
	}

	// Leave room
	h.room.Leave(roomID, clientID) // TODO: find a better way to handle unknown roomId situation
	// Delete connection
	h.connection.Delete(clientID)
	// Delete user
	h.user.Delete(clientID)

	// If user is removed than cannot inform who left the chat
	if user.ID == "<removed>" {
		return
	}

	// If user joined a chat room than get user ids in that chat
	var userIDs []string
	{
		if room, ok := h.room.Room(roomID); ok {
			userIDs = room.Users
		}
	}

	// If user joined a chat room and there is other users in chat than inform these users
	if len(userIDs) > 0 {
		res := Response{
			Body: map[string]interface{}{
				"message": "a user lost connection",
				"data":    &user,
			},
			Type: OTHER_LEFT_CHAT,
		}

		for _, userID := range userIDs {
			if c, ok := h.connection.Load(userID); ok {
				if userID == user.ID {
					continue // pass user itself
				}

				if err := c.WriteJSON(res); err != nil {
					if e := h.error(c, fiber.ErrInternalServerError); e != nil {
						h.unregister(c)
						// return
						continue
					}
				}
			}
		}
	}
}

func (h *Hub) get_rooms(conn *websocket.Conn) {
	// Load rooms
	rooms := h.room.Rooms()

	res := Response{
		Body: map[string]interface{}{
			"data": &rooms,
		},
		Type: TOPIC_ROOMS,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) change_username(req *Request) {
	// Load connection
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrInternalServerError)
		h.unregister(conn)
		return
	}

	// Read username from request body
	var username string
	{
		if tmp, ok := req.Body["username"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				username = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Set new username
	user.Username = username
	h.user.Store(user.ID, user)

	// Inform user itself here
	res := Response{
		Body: map[string]interface{}{
			"message": "your username is changed",
			"data":    &user,
		},
		Type: ME_CHANGED_USERNAME,
	}
	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	// If user joined a chat room inform users in that chat
	if room, ok := h.room.UserJoinedTo(user.ID); ok {
		res := Response{
			Body: map[string]interface{}{
				"message": "a user changed its username",
				"data":    &user,
			},
			Type: OTHER_CHANGED_USERNAME,
		}

		for _, userID := range room.Users {
			if c, ok := h.connection.Load(userID); ok {
				if userID == user.ID {
					continue // pass user itself
				}

				if err := c.WriteJSON(res); err != nil {
					if e := h.error(c, fiber.ErrInternalServerError); e != nil {
						h.unregister(c)
						// return
						continue
					}
				}
			}
		}
	}
}

func (h *Hub) join_chat(req *Request) {
	// Load connection
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrInternalServerError)
		h.unregister(conn)
		return
	}

	// Read roomId from request body
	var roomID string
	{
		if tmp, ok := req.Body["roomId"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				roomID = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Join chat room
	if ok := h.room.Join(roomID, req.ClientID); !ok {
		h.error(conn, fiber.ErrBadRequest)
		return
	}

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Load room
	room, ok := h.room.Room(roomID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Get last n messages by room
	var messages []Message
	for _, message := range h.message.GetLastN(roomID, h.Options.MaxReturnedMessage) {
		// Load owner
		owner, ok := h.user.Load(message.UserID)
		if !ok { // TODO: handle this more user friendly way
			owner.ID = "<removed>"
			owner.Username = "<removed>"
		}
		message.User = &owner
		messages = append(messages, message)
	}

	// Load online users
	var users []User
	for _, id := range h.room.Users(roomID) {
		if user, ok := h.user.Load(id); ok {
			users = append(users, user)
		}
	}

	// Inform user itself here and send room info and last n messages back
	res := Response{
		Body: map[string]interface{}{
			"message": "you joined chat",
			"data": map[string]interface{}{
				"room":     room,
				"messages": &messages,
				"users":    &users,
			},
		},
		Type: ME_JOINED_CHAT,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	// Inform users in chat
	res.Body = map[string]interface{}{
		"message": "a user joined chat",
		"data":    &user,
	}
	res.Type = OTHER_JOINED_CHAT
	for _, userID := range h.room.Users(roomID) {
		if c, ok := h.connection.Load(userID); ok {
			if userID == user.ID {
				continue // pass user itself
			}

			if err := c.WriteJSON(res); err != nil {
				if e := h.error(c, fiber.ErrInternalServerError); e != nil {
					h.unregister(c)
					// return
					continue
				}
			}
		}
	}
}

func (h *Hub) leave_chat(req *Request) {
	// Load connection
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrInternalServerError)
		h.unregister(conn)
		return
	}

	// Read roomId from request body
	var roomID string
	{
		if tmp, ok := req.Body["roomId"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				roomID = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Leave chat room
	h.room.Leave(roomID, req.ClientID)

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Inform user itself here
	res := Response{
		Body: map[string]interface{}{
			"message": "you left the chat",
		},
		Type: ME_LEFT_CHAT,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	// Inform users in chat
	res.Body = map[string]interface{}{
		"message": "a user left chat",
		"data":    &user,
	}
	res.Type = OTHER_LEFT_CHAT
	for _, userID := range h.room.Users(roomID) {
		if c, ok := h.connection.Load(userID); ok {
			if userID == user.ID {
				continue // pass user itself
			}

			if err := c.WriteJSON(res); err != nil {
				if e := h.error(c, fiber.ErrInternalServerError); e != nil {
					h.unregister(c)
					// return
					continue
				}
			}
		}
	}
}

func (h *Hub) send_message(req *Request) {
	// Load connection
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrInternalServerError)
		h.unregister(conn)
		return
	}

	// Read roomId from request body
	var roomID string
	{
		if tmp, ok := req.Body["roomId"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				roomID = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Read message from request body
	var message string
	{
		if tmp, ok := req.Body["message"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				message = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Remove old messages
	if h.message.Count(roomID) >= h.Options.MaxSavedMessage {
		h.message.Set(roomID, h.message.Get(roomID)[1:])
	}

	// Save new message
	newMessage := Message{
		ID:        req.ID,
		UserID:    req.ClientID,
		RoomID:    roomID,
		Message:   message,
		Timestamp: time.Now().Unix() * 1000, // in ms
	}
	h.message.Append(roomID, newMessage)

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Inform user itself here
	res := Response{
		Body: map[string]interface{}{
			"data": &newMessage,
		},
		Type: ME_MESSAGE_SEND,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	// Inform users in chat
	res.Type = OTHER_MESSAGE_SEND
	newMessage.User = &user
	for _, userID := range h.room.Users(roomID) {
		if c, ok := h.connection.Load(userID); ok {
			if userID == user.ID {
				continue // pass user itself
			}

			if err := c.WriteJSON(res); err != nil {
				if e := h.error(c, fiber.ErrInternalServerError); e != nil {
					h.unregister(c)
					// return
					continue
				}
			}
		}
	}
}

func (h *Hub) old_messages(req *Request) {
	// Load connection
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, fiber.ErrInternalServerError)
		h.unregister(conn)
		return
	}

	// Read roomId from request body
	var (
		roomID      string
		oldestMsgID string
	)
	{
		if tmp, ok := req.Body["roomId"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				roomID = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}

		if tmp, ok := req.Body["oldestMsgId"]; ok {
			if s, ok := tmp.(string); ok { // TODO: validate data. e.g., not empty
				oldestMsgID = s
			} else {
				h.error(conn, fiber.ErrBadRequest)
				return
			}
		} else {
			h.error(conn, fiber.ErrBadRequest)
			return
		}
	}

	// Load room
	room, ok := h.room.Room(roomID)
	if !ok {
		h.error(conn, fiber.ErrNotFound)
		return
	}

	// Get last n messages by room older than oldestMsgID
	var messages []Message
	for _, message := range h.message.GetLastN(roomID, h.Options.MaxReturnedMessage, oldestMsgID) {
		// Load owner
		owner, ok := h.user.Load(message.UserID)
		if !ok { // TODO: handle this more user friendly way
			owner.ID = "<removed>"
			owner.Username = "<removed>"
		}
		message.User = &owner
		messages = append(messages, message)
	}

	// Inform user itself here and send room info and last n messages back
	res := Response{
		Body: map[string]interface{}{
			"data": map[string]interface{}{
				"room":     room,
				"messages": &messages,
			},
		},
		Type: OLD_MESSAGES,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, fiber.ErrInternalServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) error(conn *websocket.Conn, err error) error {
	res := Response{
		Error: map[string]interface{}{
			"message": err.Error(),
		},
		Type: ERROR,
	}
	return conn.WriteJSON(res)
}
