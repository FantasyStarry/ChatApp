package controllers

import (
	"chatapp/config"
	"chatapp/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateChatRoomRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// GetChatRooms returns all chat rooms
func GetChatRooms(c *gin.Context) {
	var chatRooms []models.ChatRoom
	if err := config.DB.Preload("Creator").Find(&chatRooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat rooms"})
		return
	}

	c.JSON(http.StatusOK, chatRooms)
}

// GetChatRoom returns a specific chat room
func GetChatRoom(c *gin.Context) {
	id := c.Param("id")

	var chatRoom models.ChatRoom
	if err := config.DB.Preload("Creator").Preload("Messages.User").First(&chatRoom, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat room not found"})
		return
	}

	c.JSON(http.StatusOK, chatRoom)
}

// CreateChatRoom creates a new chat room
func CreateChatRoom(c *gin.Context) {
	var req CreateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	chatRoom := models.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID.(uint),
	}

	if err := config.DB.Create(&chatRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat room"})
		return
	}

	// Preload creator information
	config.DB.Preload("Creator").First(&chatRoom, chatRoom.ID)

	c.JSON(http.StatusCreated, chatRoom)
}

// GetChatRoomMessages returns messages for a specific chat room
func GetChatRoomMessages(c *gin.Context) {
	id := c.Param("id")

	// Parse optional query parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var messages []models.Message
	if err := config.DB.Where("chat_room_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
