package repository

import (
	"chatapp/models"

	"gorm.io/gorm"
)

// ChatRoomRepository handles chat room data operations
type ChatRoomRepository interface {
	Create(chatRoom *models.ChatRoom) error
	GetByID(id uint) (*models.ChatRoom, error)
	GetByIDWithMessages(id uint) (*models.ChatRoom, error)
	List() ([]models.ChatRoom, error)
	Update(chatRoom *models.ChatRoom) error
	Delete(id uint) error
	GetByCreatorID(creatorID uint) ([]models.ChatRoom, error)
}

type chatRoomRepository struct {
	db *gorm.DB
}

// NewChatRoomRepository creates a new chat room repository
func NewChatRoomRepository(db *gorm.DB) ChatRoomRepository {
	return &chatRoomRepository{db: db}
}

func (r *chatRoomRepository) Create(chatRoom *models.ChatRoom) error {
	return r.db.Create(chatRoom).Error
}

func (r *chatRoomRepository) GetByID(id uint) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := r.db.Preload("Creator").First(&chatRoom, id).Error
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func (r *chatRoomRepository) GetByIDWithMessages(id uint) (*models.ChatRoom, error) {
	var chatRoom models.ChatRoom
	err := r.db.Preload("Creator").Preload("Messages.User").First(&chatRoom, id).Error
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func (r *chatRoomRepository) List() ([]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom
	err := r.db.Preload("Creator").Find(&chatRooms).Error
	return chatRooms, err
}

func (r *chatRoomRepository) Update(chatRoom *models.ChatRoom) error {
	return r.db.Save(chatRoom).Error
}

func (r *chatRoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.ChatRoom{}, id).Error
}

func (r *chatRoomRepository) GetByCreatorID(creatorID uint) ([]models.ChatRoom, error) {
	var chatRooms []models.ChatRoom
	err := r.db.Where("created_by = ?", creatorID).Preload("Creator").Find(&chatRooms).Error
	return chatRooms, err
}
