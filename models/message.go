package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Content    string         `json:"content" gorm:"not null"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`
ChatRoomID uint           `json:"chat_room_id" gorm:"column:chat_room_id"`
	ChatRoom   ChatRoom       `json:"chatroom,omitempty" gorm:"foreignKey:ChatRoomID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// WebSocket message structure
type WSMessage struct {
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	ChatRoomID uint      `json:"chat_room_id"`
	Timestamp  time.Time `json:"timestamp"`
	// Authentication fields
	Token string `json:"token,omitempty"`
}
