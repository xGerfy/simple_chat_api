package main

import (
	"log"
	"simple_chat_api/internal/app"
	"simple_chat_api/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Создание приложения
	application := app.NewApp(cfg)

	// Инициализация базы данных
	if err := application.InitializeDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Инициализация маршрутов
	application.InitializeRoutes()

	// Запуск сервера
	if err := application.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
