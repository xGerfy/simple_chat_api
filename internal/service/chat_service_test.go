package service

import (
	"errors"
	"simple_chat_api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок репозитория чатов
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(chat *models.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) GetByID(id int, limit int) (*models.Chat, error) {
	args := m.Called(id, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chat), args.Error(1)
}

func (m *MockChatRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Мок репозитория сообщений
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message *models.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func TestNewChatService(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)

	service := NewChatService(mockChatRepo, mockMessageRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &chatService{}, service)
}

func TestChatService_CreateChat_Success(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	mockChatRepo.On("Create", mock.AnythingOfType("*models.Chat")).
		Return(nil).
		Run(func(args mock.Arguments) {
			chat := args.Get(0).(*models.Chat)
			chat.ID = 1
		})

	// Выполнение теста
	req := models.CreateChatRequest{Title: "Test Chat"}
	chat, err := service.CreateChat(req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, chat)
	assert.Equal(t, 1, chat.ID)
	assert.Equal(t, "Test Chat", chat.Title)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_CreateChat_EmptyTitle(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Выполнение теста
	req := models.CreateChatRequest{Title: ""}
	chat, err := service.CreateChat(req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, chat)
	assert.IsType(t, &models.ValidationError{}, err)
}

func TestChatService_CreateChat_TitleTooLong(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Выполнение теста
	req := models.CreateChatRequest{Title: string(make([]byte, 201))}
	chat, err := service.CreateChat(req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, chat)
	assert.IsType(t, &models.ValidationError{}, err)
}

func TestChatService_CreateChat_RepositoryError(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	expectedErr := errors.New("database error")
	mockChatRepo.On("Create", mock.AnythingOfType("*models.Chat")).
		Return(expectedErr)

	// Выполнение теста
	req := models.CreateChatRequest{Title: "Test Chat"}
	chat, err := service.CreateChat(req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, chat)
	assert.Equal(t, expectedErr, err)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_CreateMessage_Success(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка моков
	existingChat := &models.Chat{ID: 1, Title: "Existing Chat"}
	mockChatRepo.On("GetByID", 1, 1).Return(existingChat, nil)

	mockMessageRepo.On("Create", mock.AnythingOfType("*models.Message")).
		Return(nil).
		Run(func(args mock.Arguments) {
			msg := args.Get(0).(*models.Message)
			msg.ID = 1
		})

	// Выполнение теста
	req := models.CreateMessageRequest{Text: "Hello World"}
	message, err := service.CreateMessage(1, req)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, 1, message.ChatID)
	assert.Equal(t, "Hello World", message.Text)
	mockChatRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_CreateMessage_ChatNotFound(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	mockChatRepo.On("GetByID", 999, 1).Return(nil, nil)

	// Выполнение теста
	req := models.CreateMessageRequest{Text: "Hello World"}
	message, err := service.CreateMessage(999, req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.IsType(t, &NotFoundError{}, err)
	assert.Equal(t, "chat not found", err.Error())
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_CreateMessage_EmptyText(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Выполнение теста
	req := models.CreateMessageRequest{Text: ""}
	message, err := service.CreateMessage(1, req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.IsType(t, &models.ValidationError{}, err)

	// Убеждаемся, что GetByID НЕ вызывался
	mockChatRepo.AssertNotCalled(t, "GetByID")
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_CreateMessage_TextTooLong(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Выполнение теста
	req := models.CreateMessageRequest{Text: string(make([]byte, 5001))}
	message, err := service.CreateMessage(1, req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.IsType(t, &models.ValidationError{}, err)

	// Убеждаемся, что GetByID НЕ вызывался
	mockChatRepo.AssertNotCalled(t, "GetByID")
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_GetChatWithMessages_Success(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	expectedChat := &models.Chat{
		ID:    1,
		Title: "Test Chat",
		Messages: []models.Message{
			{ID: 1, ChatID: 1, Text: "Message 1"},
			{ID: 2, ChatID: 1, Text: "Message 2"},
		},
	}
	mockChatRepo.On("GetByID", 1, 20).Return(expectedChat, nil)

	// Выполнение теста
	chat, err := service.GetChatWithMessages(1, 20)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, chat)
	assert.Equal(t, 1, chat.ID)
	assert.Equal(t, "Test Chat", chat.Title)
	assert.Len(t, chat.Messages, 2)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_GetChatWithMessages_NotFound(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	mockChatRepo.On("GetByID", 999, 20).Return(nil, nil)

	// Выполнение теста
	chat, err := service.GetChatWithMessages(999, 20)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, chat)
	assert.IsType(t, &NotFoundError{}, err)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_GetChatWithMessages_LimitExceeded(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	expectedChat := &models.Chat{
		ID:    1,
		Title: "Test Chat",
	}
	mockChatRepo.On("GetByID", 1, 100).Return(expectedChat, nil)

	// Выполнение теста (лимит больше максимального)
	chat, err := service.GetChatWithMessages(1, 150)

	// Проверки
	assert.NoError(t, err)
	assert.NotNil(t, chat)
	// Проверяем, что вызвался с лимитом 100
	mockChatRepo.AssertCalled(t, "GetByID", 1, 100)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_DeleteChat_Success(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	mockChatRepo.On("Delete", 1).Return(nil)

	// Выполнение теста
	err := service.DeleteChat(1)

	// Проверки
	assert.NoError(t, err)
	mockChatRepo.AssertExpectations(t)
}

func TestChatService_DeleteChat_RepositoryError(t *testing.T) {
	mockChatRepo := new(MockChatRepository)
	mockMessageRepo := new(MockMessageRepository)
	service := NewChatService(mockChatRepo, mockMessageRepo)

	// Настройка мока
	expectedErr := errors.New("database error")
	mockChatRepo.On("Delete", 1).Return(expectedErr)

	// Выполнение теста
	err := service.DeleteChat(1)

	// Проверки
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockChatRepo.AssertExpectations(t)
}
