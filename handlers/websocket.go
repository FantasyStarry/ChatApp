package handlers

import (
	"chatapp/config"
	"chatapp/models"
	"chatapp/service"
	"chatapp/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// InitWebSocketUpgrader initializes the WebSocket upgrader with config values
func InitWebSocketUpgrader() {
	if config.GlobalConfig != nil {
		upgrader.ReadBufferSize = config.GlobalConfig.WebSocket.ReadBufferSize
		upgrader.WriteBufferSize = config.GlobalConfig.WebSocket.WriteBufferSize
	}
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Chat room specific clients
	chatRooms map[uint]map[*Client]bool

	// Message service for database operations
	messageService service.MessageService
}

type Client struct {
	hub             *Hub
	conn            *websocket.Conn
	send            chan []byte
	userID          uint
	username        string
	chatRoomID      uint
	isAuthenticated bool
}

func NewHub(messageService service.MessageService) *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		chatRooms:      make(map[uint]map[*Client]bool),
		messageService: messageService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// Add client to specific chat room
			if h.chatRooms[client.chatRoomID] == nil {
				h.chatRooms[client.chatRoomID] = make(map[*Client]bool)
			}
			h.chatRooms[client.chatRoomID][client] = true

			log.Printf("Client %s joined chat room %d", client.username, client.chatRoomID)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.chatRooms[client.chatRoomID], client)
				close(client.send)
				log.Printf("Client %s left chat room %d", client.username, client.chatRoomID)
			}

		case message := <-h.broadcast:
			// Broadcast to all clients (this could be modified to be room-specific)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(chatRoomID uint, message []byte) {
	if clients, ok := h.chatRooms[chatRoomID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(h.chatRooms[chatRoomID], client)
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var wsMsg models.WSMessage
		err := c.conn.ReadJSON(&wsMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Handle authentication message
		if wsMsg.Type == "auth" {
			if err := c.handleAuthMessage(wsMsg); err != nil {
				log.Printf("Authentication failed for client: %v", err)
				c.conn.Close()
				return
			}
			continue // Don't broadcast or save auth messages
		}

		// Check if client is authenticated for non-auth messages
		if !c.isAuthenticated {
			log.Printf("Unauthenticated client tried to send message")
			c.conn.Close()
			return
		}

		var message *models.Message
		var saveErr error

		// Save message to database using service layer based on message type
		if wsMsg.Type == "file" {
			message, saveErr = c.hub.messageService.CreateFileMessage(wsMsg.Content, c.userID, c.chatRoomID)
		} else {
			message, saveErr = c.hub.messageService.CreateMessage(wsMsg.Content, c.userID, c.chatRoomID)
		}
		
		if saveErr != nil {
			log.Printf("Failed to save message: %v", saveErr)
			continue
		}

		// Create response message
		responseMsg := models.WSMessage{
			Type:       message.Type,
			Content:    message.Content,
			UserID:     message.UserID,
			Username:   message.User.Username,
			ChatRoomID: message.ChatRoomID,
			Timestamp:  message.CreatedAt,
		}

		// Broadcast to all clients in the same chat room
		if msgBytes, err := json.Marshal(responseMsg); err == nil {
			c.hub.BroadcastToRoom(c.chatRoomID, msgBytes)
		}
	}
}

// handleAuthMessage processes authentication messages from WebSocket clients
func (c *Client) handleAuthMessage(wsMsg models.WSMessage) error {
	// Validate the token
	claims, err := utils.ValidateToken(wsMsg.Token)
	if err != nil {
		return err
	}

	// Set client authentication details
	c.userID = claims.UserID
	c.username = claims.Username
	// c.chatRoomID = wsMsg.ChatRoomID
	c.isAuthenticated = true

	log.Printf("Client authenticated: user_id=%d, username=%s, chatroom_id=%d",
		c.userID, c.username, c.chatRoomID)

	// Register the client with the hub after successful authentication
	c.hub.register <- c

	// Send authentication success response
	response := models.WSMessage{
		Type:      "auth_success",
		Content:   "Authentication successful",
		Timestamp: time.Now(),
	}

	if msgBytes, err := json.Marshal(response); err == nil {
		select {
		case c.send <- msgBytes:
		default:
			close(c.send)
		}
	}

	return nil
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// GlobalHub will be initialized in main.go with proper dependencies
var GlobalHub *Hub

// InitializeHub initializes the global hub with message service
func InitializeHub(messageService service.MessageService) {
	GlobalHub = NewHub(messageService)
}

func HandleWebSocket(c *gin.Context) {
	chatRoomIDStr := c.Param("chatroom_id")
	chatRoomID, err := strconv.ParseUint(chatRoomIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid chat room ID")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:             GlobalHub,
		conn:            conn,
		send:            make(chan []byte, 256),
		userID:          0,  // Will be set during authentication
		username:        "", // Will be set during authentication
		chatRoomID:      uint(chatRoomID),
		isAuthenticated: false,
	}

	// Note: We don't register the client immediately anymore
	// Registration will happen after successful authentication

	go client.writePump()
	go client.readPump()
}
