package repository

import (
	"errors"
	"simple_chat_api/internal/models"

	"gorm.io/gorm"
)

type ChatRepository interface {
	Create(chat *models.Chat) error
	GetByID(id int, limit int) (*models.Chat, error)
	Delete(id int) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(chat *models.Chat) error {
	return r.db.Create(chat).Error
}

func (r *chatRepository) GetByID(id int, limit int) (*models.Chat, error) {
	var chat models.Chat

	// Загружаем чат
	err := r.db.First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Загружаем сообщения с лимитом
	err = r.db.Model(&chat).
		Limit(limit).
		Order("created_at DESC").
		Association("Messages").
		Find(&chat.Messages)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *chatRepository) Delete(id int) error {
	return r.db.Delete(&models.Chat{}, id).Error
}
