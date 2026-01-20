package models

import (
	"strings"
	"time"
)

type Chat struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Messages  []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE;" json:"messages,omitempty"`
}

type CreateChatRequest struct {
	Title string `json:"title" binding:"required"`
}

func (r *CreateChatRequest) Validate() error {
	title := strings.TrimSpace(r.Title)

	if title == "" {
		return &ValidationError{Field: "title", Message: "title cannot be empty"}
	}

	if len(title) > 200 {
		return &ValidationError{Field: "title", Message: "title must be less than 200 characters"}
	}

	r.Title = title
	return nil
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
