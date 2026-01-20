package app

import (
	"log"
	"net/http"
	"simple_chat_api/internal/config"
	"simple_chat_api/internal/handlers"
	"simple_chat_api/internal/repository"
	"simple_chat_api/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	db     *gorm.DB
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) InitializeDB() error {
	dsn := "host=" + a.config.DBHost +
		" port=" + a.config.DBPort +
		" user=" + a.config.DBUser +
		" password=" + a.config.DBPassword +
		" dbname=" + a.config.DBName +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	a.db = db
	log.Println("Database connection established")

	return nil
}

func (a *App) InitializeRoutes() {
	// Инициализация репозиториев
	chatRepo := repository.NewChatRepository(a.db)
	messageRepo := repository.NewMessageRepository(a.db)

	// Инициализация сервиса
	chatService := service.NewChatService(chatRepo, messageRepo)

	// Инициализация обработчиков
	chatHandler := handlers.NewChatHandler(chatService)

	// Настройка маршрутов
	mux := http.NewServeMux()

	mux.HandleFunc("POST /chats/", chatHandler.CreateChat)
	mux.HandleFunc("POST /chats/{id}/messages/", chatHandler.CreateMessage)
	mux.HandleFunc("GET /chats/{id}", chatHandler.GetChat)
	mux.HandleFunc("DELETE /chats/{id}", chatHandler.DeleteChat)

	a.server = &http.Server{
		Addr:    ":" + a.config.ServerPort,
		Handler: mux,
	}
}

func (a *App) Run() error {
	log.Printf("Server starting on port %s", a.config.ServerPort)
	return a.server.ListenAndServe()
}
