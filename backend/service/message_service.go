package service

import (
	"chatapp/models"
	"chatapp/repository"
	"errors"
)

// MessageService handles message business logic
type MessageService interface {
	CreateMessage(content string, userID, chatRoomID uint) (*models.Message, error)
	CreateFileMessage(content string, userID, chatRoomID uint) (*models.Message, error)
	GetMessage(id uint) (*models.Message, error)
	GetChatRoomMessages(chatRoomID uint, limit, offset int) ([]models.Message, error)
	GetUserMessages(userID uint, limit, offset int) ([]models.Message, error)
	UpdateMessage(id uint, content string, userID uint) (*models.Message, error)
	DeleteMessage(id uint, userID uint) error
	GetRecentMessages(chatRoomID uint, limit int) ([]models.Message, error)
	GetMessageCount(chatRoomID uint) (int64, error)
}

type messageService struct {
	messageRepo  repository.MessageRepository
	userRepo     repository.UserRepository
	chatRoomRepo repository.ChatRoomRepository
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo repository.MessageRepository, userRepo repository.UserRepository, chatRoomRepo repository.ChatRoomRepository) MessageService {
	return &messageService{
		messageRepo:  messageRepo,
		userRepo:     userRepo,
		chatRoomRepo: chatRoomRepo,
	}
}

func (s *messageService) CreateMessage(content string, userID, chatRoomID uint) (*models.Message, error) {
	// Validate content
	if content == "" {
		return nil, errors.New("message content is required")
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate chat room exists
	_, err = s.chatRoomRepo.GetByID(chatRoomID)
	if err != nil {
		return nil, errors.New("chat room not found")
	}

	message := &models.Message{
		Content:    content,
		UserID:     userID,
		ChatRoomID: chatRoomID,
		Type:       "message", // Default type is message
	}

	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, errors.New("failed to create message")
	}

	// Return message with user and chat room information
	return s.messageRepo.GetByID(message.ID)
}

func (s *messageService) CreateFileMessage(content string, userID, chatRoomID uint) (*models.Message, error) {
	// Validate content
	if content == "" {
		return nil, errors.New("file information is required")
	}

	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate chat room exists
	_, err = s.chatRoomRepo.GetByID(chatRoomID)
	if err != nil {
		return nil, errors.New("chat room not found")
	}

	message := &models.Message{
		Content:    content,
		UserID:     userID,
		ChatRoomID: chatRoomID,
		Type:       "file",
	}

	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, errors.New("failed to create file message")
	}

	// Return message with user and chat room information
	return s.messageRepo.GetByID(message.ID)
}

func (s *messageService) GetMessage(id uint) (*models.Message, error) {
	message, err := s.messageRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("message not found")
	}
	return message, nil
}

func (s *messageService) GetChatRoomMessages(chatRoomID uint, limit, offset int) ([]models.Message, error) {
	// Validate chat room exists
	_, err := s.chatRoomRepo.GetByID(chatRoomID)
	if err != nil {
		return nil, errors.New("chat room not found")
	}

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}

	return s.messageRepo.GetByChatRoomID(chatRoomID, limit, offset)
}

func (s *messageService) GetUserMessages(userID uint, limit, offset int) ([]models.Message, error) {
	// Validate user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}

	return s.messageRepo.GetByUserID(userID, limit, offset)
}

func (s *messageService) UpdateMessage(id uint, content string, userID uint) (*models.Message, error) {
	// Validate content
	if content == "" {
		return nil, errors.New("message content is required")
	}

	// Get existing message
	message, err := s.messageRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("message not found")
	}

	// Check if user is the author
	if message.UserID != userID {
		return nil, errors.New("only the author can update this message")
	}

	message.Content = content
	err = s.messageRepo.Update(message)
	if err != nil {
		return nil, errors.New("failed to update message")
	}

	return message, nil
}

func (s *messageService) DeleteMessage(id uint, userID uint) error {
	// Get existing message
	message, err := s.messageRepo.GetByID(id)
	if err != nil {
		return errors.New("message not found")
	}

	// Check if user is the author
	if message.UserID != userID {
		return errors.New("only the author can delete this message")
	}

	return s.messageRepo.Delete(id)
}

func (s *messageService) GetRecentMessages(chatRoomID uint, limit int) ([]models.Message, error) {
	// Validate chat room exists
	_, err := s.chatRoomRepo.GetByID(chatRoomID)
	if err != nil {
		return nil, errors.New("chat room not found")
	}

	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	return s.messageRepo.GetRecentMessages(chatRoomID, limit)
}

func (s *messageService) GetMessageCount(chatRoomID uint) (int64, error) {
	// Validate chat room exists
	_, err := s.chatRoomRepo.GetByID(chatRoomID)
	if err != nil {
		return 0, errors.New("chat room not found")
	}

	return s.messageRepo.CountByChatRoomID(chatRoomID)
}
