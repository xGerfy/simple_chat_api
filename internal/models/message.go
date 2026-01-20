package models

import (
	"strings"
	"time"
)

type Message struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    int       `gorm:"not null;index" json:"chat_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CreateMessageRequest struct {
	Text string `json:"text" binding:"required"`
}

func (r *CreateMessageRequest) Validate() error {
	text := strings.TrimSpace(r.Text)

	if text == "" {
		return &ValidationError{Field: "text", Message: "text cannot be empty"}
	}

	if len(text) > 5000 {
		return &ValidationError{Field: "text", Message: "text must be less than 5000 characters"}
	}

	r.Text = text
	return nil
}
