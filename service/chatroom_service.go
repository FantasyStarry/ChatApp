package service

import (
	"chatapp/models"
	"chatapp/repository"
	"errors"
)

// ChatRoomService handles chat room business logic
type ChatRoomService interface {
	CreateChatRoom(name, description string, creatorID uint) (*models.ChatRoom, error)
	GetChatRoom(id uint) (*models.ChatRoom, error)
	GetChatRoomWithMessages(id uint) (*models.ChatRoom, error)
	GetAllChatRooms() ([]models.ChatRoom, error)
	UpdateChatRoom(id uint, name, description string, userID uint) (*models.ChatRoom, error)
	DeleteChatRoom(id uint, userID uint) error
	GetUserChatRooms(userID uint) ([]models.ChatRoom, error)
}

type chatRoomService struct {
	chatRoomRepo repository.ChatRoomRepository
	userRepo     repository.UserRepository
}

// NewChatRoomService creates a new chat room service
func NewChatRoomService(chatRoomRepo repository.ChatRoomRepository, userRepo repository.UserRepository) ChatRoomService {
	return &chatRoomService{
		chatRoomRepo: chatRoomRepo,
		userRepo:     userRepo,
	}
}

func (s *chatRoomService) CreateChatRoom(name, description string, creatorID uint) (*models.ChatRoom, error) {
	// Validate creator exists
	_, err := s.userRepo.GetByID(creatorID)
	if err != nil {
		return nil, errors.New("creator not found")
	}

	// Validate name
	if name == "" {
		return nil, errors.New("chat room name is required")
	}

	chatRoom := &models.ChatRoom{
		Name:        name,
		Description: description,
		CreatedBy:   creatorID,
	}

	err = s.chatRoomRepo.Create(chatRoom)
	if err != nil {
		return nil, errors.New("failed to create chat room")
	}

	// Return chat room with creator information
	return s.chatRoomRepo.GetByID(chatRoom.ID)
}

func (s *chatRoomService) GetChatRoom(id uint) (*models.ChatRoom, error) {
	chatRoom, err := s.chatRoomRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("chat room not found")
	}
	return chatRoom, nil
}

func (s *chatRoomService) GetChatRoomWithMessages(id uint) (*models.ChatRoom, error) {
	chatRoom, err := s.chatRoomRepo.GetByIDWithMessages(id)
	if err != nil {
		return nil, errors.New("chat room not found")
	}
	return chatRoom, nil
}

func (s *chatRoomService) GetAllChatRooms() ([]models.ChatRoom, error) {
	return s.chatRoomRepo.List()
}

func (s *chatRoomService) UpdateChatRoom(id uint, name, description string, userID uint) (*models.ChatRoom, error) {
	// Get existing chat room
	chatRoom, err := s.chatRoomRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("chat room not found")
	}

	// Check if user is the creator
	if chatRoom.CreatedBy != userID {
		return nil, errors.New("only the creator can update this chat room")
	}

	// Update fields
	if name != "" {
		chatRoom.Name = name
	}
	chatRoom.Description = description

	err = s.chatRoomRepo.Update(chatRoom)
	if err != nil {
		return nil, errors.New("failed to update chat room")
	}

	return chatRoom, nil
}

func (s *chatRoomService) DeleteChatRoom(id uint, userID uint) error {
	// Get existing chat room
	chatRoom, err := s.chatRoomRepo.GetByID(id)
	if err != nil {
		return errors.New("chat room not found")
	}

	// Check if user is the creator
	if chatRoom.CreatedBy != userID {
		return errors.New("only the creator can delete this chat room")
	}

	return s.chatRoomRepo.Delete(id)
}

func (s *chatRoomService) GetUserChatRooms(userID uint) ([]models.ChatRoom, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.chatRoomRepo.GetByCreatorID(userID)
}
