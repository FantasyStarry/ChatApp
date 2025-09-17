package repository

import (
	"chatapp/models"

	"gorm.io/gorm"
)

// MessageRepository handles message data operations
type MessageRepository interface {
	Create(message *models.Message) error
	GetByID(id uint) (*models.Message, error)
	GetByChatRoomID(chatRoomID uint, limit, offset int) ([]models.Message, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Message, error)
	Update(message *models.Message) error
	Delete(id uint) error
	GetRecentMessages(chatRoomID uint, limit int) ([]models.Message, error)
	CountByChatRoomID(chatRoomID uint) (int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetByID(id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("User").Preload("ChatRoom").First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetByChatRoomID(chatRoomID uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("chat_room_id = ?", chatRoomID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) GetByUserID(userID uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("user_id = ?", userID).
		Preload("User").
		Preload("ChatRoom").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *messageRepository) Delete(id uint) error {
	return r.db.Delete(&models.Message{}, id).Error
}

func (r *messageRepository) GetRecentMessages(chatRoomID uint, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("chat_room_id = ?", chatRoomID).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) CountByChatRoomID(chatRoomID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).Where("chat_room_id = ?", chatRoomID).Count(&count).Error
	return count, err
}
