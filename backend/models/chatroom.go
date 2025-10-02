package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatRoom struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedBy   uint           `json:"created_by"`
	Creator     User           `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Messages    []Message      `json:"messages,omitempty" gorm:"foreignKey:ChatRoomID"`
}
