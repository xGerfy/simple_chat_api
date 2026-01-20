package service

import (
	"simple_chat_api/internal/models"
	"simple_chat_api/internal/repository"
)

type ChatService interface {
	CreateChat(req models.CreateChatRequest) (*models.Chat, error)
	CreateMessage(chatID int, req models.CreateMessageRequest) (*models.Message, error)
	GetChatWithMessages(id int, limit int) (*models.Chat, error)
	DeleteChat(id int) error
}

type chatService struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
}

func NewChatService(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository) ChatService {
	return &chatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

func (s *chatService) CreateChat(req models.CreateChatRequest) (*models.Chat, error) {
	// Валидация
	if err := req.Validate(); err != nil {
		return nil, err
	}

	chat := &models.Chat{
		Title: req.Title,
	}

	err := s.chatRepo.Create(chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *chatService) CreateMessage(chatID int, req models.CreateMessageRequest) (*models.Message, error) {
	// Валидация
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Проверяем существование чата
	chat, err := s.chatRepo.GetByID(chatID, 1)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, &NotFoundError{Resource: "chat", ID: chatID}
	}

	message := &models.Message{
		ChatID: chatID,
		Text:   req.Text,
	}

	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *chatService) GetChatWithMessages(id int, limit int) (*models.Chat, error) {
	if limit > 100 {
		limit = 100
	}

	chat, err := s.chatRepo.GetByID(id, limit)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, &NotFoundError{Resource: "chat", ID: id}
	}

	return chat, nil
}

func (s *chatService) DeleteChat(id int) error {
	return s.chatRepo.Delete(id)
}

// Ошибки
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return e.Resource + " not found"
}
