package repository

import (
	"simple_chat_api/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMessageMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestMessageRepository_Create_Success(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	message := &models.Message{
		ChatID: 1,
		Text:   "Test message",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "Test message", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message)

	assert.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, 1, message.ChatID)
	assert.Equal(t, "Test message", message.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_Error(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	message := &models.Message{
		ChatID: 1,
		Text:   "Test message",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "Test message", sqlmock.AnyArg()).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Create(message)

	assert.Error(t, err)
	assert.Equal(t, 0, message.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_MultipleMessages(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	// Тест 1: Первое сообщение
	message1 := &models.Message{
		ChatID: 1,
		Text:   "First message",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "First message", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message1)
	assert.NoError(t, err)
	assert.Equal(t, 1, message1.ID)

	// Сбрасываем мок для следующего теста
	mock.ExpectationsWereMet()

	// Тест 2: Второе сообщение
	message2 := &models.Message{
		ChatID: 1,
		Text:   "Second message",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "Second message", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	err = repo.Create(message2)
	assert.NoError(t, err)
	assert.Equal(t, 2, message2.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_DifferentChats(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	// Сообщение для первого чата
	message1 := &models.Message{
		ChatID: 1,
		Text:   "Message for chat 1",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "Message for chat 1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message1)
	assert.NoError(t, err)

	// Сбрасываем мок
	mock.ExpectationsWereMet()

	// Сообщение для второго чата
	message2 := &models.Message{
		ChatID: 2,
		Text:   "Message for chat 2",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(2, "Message for chat 2", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	err = repo.Create(message2)
	assert.NoError(t, err)
	assert.Equal(t, 2, message2.ChatID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_LongText(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	longText := "This is a very long message. " +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris."

	message := &models.Message{
		ChatID: 1,
		Text:   longText,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, longText, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message)

	assert.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, longText, message.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_WithSpecificTime(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	specificTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	message := &models.Message{
		ChatID:    1,
		Text:      "Message with specific time",
		CreatedAt: specificTime,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "Message with specific time", specificTime).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message)

	assert.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, specificTime, message.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Create_EmptyText(t *testing.T) {
	db, mock := setupMessageMockDB(t)
	repo := NewMessageRepository(db)

	// Репозиторий не должен валидировать данные, только вставлять
	// Валидация должна быть на уровне сервиса
	message := &models.Message{
		ChatID: 1,
		Text:   "", // Пустой текст
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "messages" ("chat_id","text","created_at") VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, "", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(message)

	// Репозиторий не валидирует, поэтому ошибки быть не должно
	assert.NoError(t, err)
	assert.Equal(t, 1, message.ID)
	assert.Equal(t, "", message.Text)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewMessageRepository(t *testing.T) {
	db, mock := setupMessageMockDB(t)

	repo := NewMessageRepository(db)

	assert.NotNil(t, repo)
	assert.IsType(t, &messageRepository{}, repo)

	// Проверяем, что репозиторий реализует интерфейс
	var repoInterface MessageRepository = repo
	assert.NotNil(t, repoInterface)

	// Убедимся, что мок не использовался
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_Interface(t *testing.T) {
	// Проверяем, что структура реализует интерфейс
	var _ MessageRepository = &messageRepository{}

	// Создаем реальный экземпляр для проверки
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := NewMessageRepository(gormDB)

	// Проверяем тип
	assert.IsType(t, &messageRepository{}, repo)

	// Проверяем, что можно присвоить интерфейсу
	var messageRepo MessageRepository = repo
	assert.NotNil(t, messageRepo)

	assert.NoError(t, mock.ExpectationsWereMet())
}
