package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"simple_chat_api/internal/models"
	"simple_chat_api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок сервиса
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(req models.CreateChatRequest) (*models.Chat, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chat), args.Error(1)
}

func (m *MockChatService) CreateMessage(chatID int, req models.CreateMessageRequest) (*models.Message, error) {
	args := m.Called(chatID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockChatService) GetChatWithMessages(id int, limit int) (*models.Chat, error) {
	args := m.Called(id, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chat), args.Error(1)
}

func (m *MockChatService) DeleteChat(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateChatHandler_Success(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &models.Chat{
		ID:    1,
		Title: "New Chat",
	}

	mockService.On("CreateChat", models.CreateChatRequest{Title: "New Chat"}).
		Return(expectedChat, nil)

	// Выполнение
	reqBody := `{"title": "New Chat"}`
	req := httptest.NewRequest("POST", "/chats/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "New Chat", response.Title)

	mockService.AssertExpectations(t)
}

func TestCreateChatHandler_ValidationError(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	validationErr := &models.ValidationError{
		Field:   "title",
		Message: "title cannot be empty",
	}

	mockService.On("CreateChat", models.CreateChatRequest{Title: ""}).
		Return(nil, validationErr)

	// Выполнение
	reqBody := `{"title": ""}`
	req := httptest.NewRequest("POST", "/chats/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "title cannot be empty", response["error"])

	mockService.AssertExpectations(t)
}

func TestCreateChatHandler_InvalidJSON(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	// Выполнение
	reqBody := `{"title": }`
	req := httptest.NewRequest("POST", "/chats/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateChatHandler_InternalServerError(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("CreateChat", models.CreateChatRequest{Title: "New Chat"}).
		Return(nil, errors.New("database error"))

	// Выполнение
	reqBody := `{"title": "New Chat"}`
	req := httptest.NewRequest("POST", "/chats/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}

func TestCreateMessageHandler_Success(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedMessage := &models.Message{
		ID:     1,
		ChatID: 1,
		Text:   "Hello World",
	}

	mockService.On("CreateMessage", 1, models.CreateMessageRequest{Text: "Hello World"}).
		Return(expectedMessage, nil)

	// Выполнение
	reqBody := `{"text": "Hello World"}`
	req := httptest.NewRequest("POST", "/chats/1/messages/", bytes.NewBufferString(reqBody))
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateMessage(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Message
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, 1, response.ChatID)
	assert.Equal(t, "Hello World", response.Text)

	mockService.AssertExpectations(t)
}

func TestCreateMessageHandler_ChatNotFound(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	notFoundErr := &service.NotFoundError{
		Resource: "chat",
		ID:       999,
	}

	mockService.On("CreateMessage", 999, models.CreateMessageRequest{Text: "Hello World"}).
		Return(nil, notFoundErr)

	// Выполнение
	reqBody := `{"text": "Hello World"}`
	req := httptest.NewRequest("POST", "/chats/999/messages/", bytes.NewBufferString(reqBody))
	req.SetPathValue("id", "999")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateMessage(rr, req)

	// Проверки
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "chat not found")

	mockService.AssertExpectations(t)
}

func TestCreateMessageHandler_InvalidChatID(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	// Выполнение
	reqBody := `{"text": "Hello World"}`
	req := httptest.NewRequest("POST", "/chats/invalid/messages/", bytes.NewBufferString(reqBody))
	req.SetPathValue("id", "invalid")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateMessage(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid chat ID")
}

func TestGetChatHandler_Success(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &models.Chat{
		ID:    1,
		Title: "Test Chat",
		Messages: []models.Message{
			{ID: 1, ChatID: 1, Text: "Message 1"},
			{ID: 2, ChatID: 1, Text: "Message 2"},
		},
	}

	mockService.On("GetChatWithMessages", 1, 20).Return(expectedChat, nil)

	// Выполнение
	req := httptest.NewRequest("GET", "/chats/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler.GetChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Chat
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test Chat", response.Title)
	assert.Len(t, response.Messages, 2)

	mockService.AssertExpectations(t)
}

func TestGetChatHandler_WithCustomLimit(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &models.Chat{
		ID:    1,
		Title: "Test Chat",
	}

	mockService.On("GetChatWithMessages", 1, 50).Return(expectedChat, nil)

	// Выполнение
	req := httptest.NewRequest("GET", "/chats/1?limit=50", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler.GetChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetChatHandler_InvalidLimit(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	expectedChat := &models.Chat{
		ID:    1,
		Title: "Test Chat",
	}

	// При невалидном лимите должен использоваться дефолтный (20)
	mockService.On("GetChatWithMessages", 1, 20).Return(expectedChat, nil)

	// Выполнение (невалидный лимит)
	req := httptest.NewRequest("GET", "/chats/1?limit=invalid", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler.GetChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetChatHandler_ChatNotFound(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	notFoundErr := &service.NotFoundError{
		Resource: "chat",
		ID:       999,
	}

	mockService.On("GetChatWithMessages", 999, 20).Return(nil, notFoundErr)

	// Выполнение
	req := httptest.NewRequest("GET", "/chats/999", nil)
	req.SetPathValue("id", "999")

	rr := httptest.NewRecorder()
	handler.GetChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "chat not found")

	mockService.AssertExpectations(t)
}

func TestDeleteChatHandler_Success(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("DeleteChat", 1).Return(nil)

	// Выполнение
	req := httptest.NewRequest("DELETE", "/chats/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler.DeleteChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestDeleteChatHandler_InvalidChatID(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	// Выполнение
	req := httptest.NewRequest("DELETE", "/chats/invalid", nil)
	req.SetPathValue("id", "invalid")

	rr := httptest.NewRecorder()
	handler.DeleteChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid chat ID")
}

func TestDeleteChatHandler_InternalServerError(t *testing.T) {
	// Подготовка
	mockService := new(MockChatService)
	handler := NewChatHandler(mockService)

	mockService.On("DeleteChat", 1).Return(errors.New("database error"))

	// Выполнение
	req := httptest.NewRequest("DELETE", "/chats/1", nil)
	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()
	handler.DeleteChat(rr, req)

	// Проверки
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	mockService.AssertExpectations(t)
}
