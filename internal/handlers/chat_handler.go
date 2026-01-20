package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"simple_chat_api/internal/models"
	"simple_chat_api/internal/service"
	"strconv"
)

type ChatHandler struct {
	service service.ChatService
}

func NewChatHandler(service service.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req models.CreateChatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chat, err := h.service.CreateChat(req)
	if err != nil {
		if _, ok := err.(*models.ValidationError); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		log.Printf("Error creating chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var req models.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message, err := h.service.CreateMessage(chatID, req)
	if err != nil {
		if _, ok := err.(*service.NotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if _, ok := err.(*models.ValidationError); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		log.Printf("Error creating message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	limit := 20
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 20
		}
	}

	chat, err := h.service.GetChatWithMessages(chatID, limit)
	if err != nil {
		if _, ok := err.(*service.NotFoundError); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Error getting chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	chatIDStr := r.PathValue("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteChat(chatID)
	if err != nil {
		log.Printf("Error deleting chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
