package controllers

import (
	"chatapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatRoomController struct {
	chatRoomService service.ChatRoomService
	messageService  service.MessageService
}

// NewChatRoomController creates a new chat room controller
func NewChatRoomController(chatRoomService service.ChatRoomService, messageService service.MessageService) *ChatRoomController {
	return &ChatRoomController{
		chatRoomService: chatRoomService,
		messageService:  messageService,
	}
}

type CreateChatRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// GetChatRooms returns all chat rooms
func (ctrl *ChatRoomController) GetChatRooms(c *gin.Context) {
	chatRooms, err := ctrl.chatRoomService.GetAllChatRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatRooms)
}

// GetChatRoom returns a specific chat room
func (ctrl *ChatRoomController) GetChatRoom(c *gin.Context) {
	id := c.Param("id")
	chatRoomID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
		return
	}

	chatRoom, err := ctrl.chatRoomService.GetChatRoomWithMessages(uint(chatRoomID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatRoom)
}

// CreateChatRoom creates a new chat room
func (ctrl *ChatRoomController) CreateChatRoom(c *gin.Context) {
	var req CreateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	chatRoom, err := ctrl.chatRoomService.CreateChatRoom(req.Name, req.Description, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chatRoom)
}

// GetChatRoomMessages returns messages for a specific chat room
func (ctrl *ChatRoomController) GetChatRoomMessages(c *gin.Context) {
	id := c.Param("id")
	chatRoomID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat room ID"})
		return
	}

	// Parse optional query parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	messages, err := ctrl.messageService.GetChatRoomMessages(uint(chatRoomID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
