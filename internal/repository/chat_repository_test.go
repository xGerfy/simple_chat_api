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

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	// Используем QueryMatcherEqual для точного сравнения
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

func TestChatRepository_Create_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	chat := &models.Chat{
		Title: "Test Chat",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "chats" ("title","created_at") VALUES ($1,$2) RETURNING "id"`).
		WithArgs("Test Chat", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(chat)

	assert.NoError(t, err)
	assert.Equal(t, 1, chat.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_Create_Error(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	chat := &models.Chat{
		Title: "Test Chat",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "chats" ("title","created_at") VALUES ($1,$2) RETURNING "id"`).
		WithArgs("Test Chat", sqlmock.AnyArg()).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Create(chat)

	assert.Error(t, err)
	assert.Equal(t, 0, chat.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	createdAt := time.Now()

	// Запрос для чата
	chatRows := sqlmock.NewRows([]string{"id", "title", "created_at"}).
		AddRow(1, "Test Chat", createdAt)

	mock.ExpectQuery(`SELECT * FROM "chats" WHERE "chats"."id" = $1 ORDER BY "chats"."id" LIMIT $2`).
		WithArgs(1, 1).
		WillReturnRows(chatRows)

	// Запрос для сообщений С LIMIT (второй параметр)
	messageRows := sqlmock.NewRows([]string{"id", "chat_id", "text", "created_at"}).
		AddRow(1, 1, "Message 1", createdAt.Add(time.Minute)).
		AddRow(2, 1, "Message 2", createdAt.Add(2*time.Minute))

	// ИСПРАВЛЕНО: Добавили LIMIT $2
	mock.ExpectQuery(`SELECT * FROM "messages" WHERE "messages"."chat_id" = $1 ORDER BY created_at DESC LIMIT $2`).
		WithArgs(1, 20). // Второй аргумент - лимит
		WillReturnRows(messageRows)

	chat, err := repo.GetByID(1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, chat)
	assert.Equal(t, 1, chat.ID)
	assert.Equal(t, "Test Chat", chat.Title)
	assert.Len(t, chat.Messages, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetByID_WithLimit(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	createdAt := time.Now()

	chatRows := sqlmock.NewRows([]string{"id", "title", "created_at"}).
		AddRow(1, "Test Chat", createdAt)

	mock.ExpectQuery(`SELECT * FROM "chats" WHERE "chats"."id" = $1 ORDER BY "chats"."id" LIMIT $2`).
		WithArgs(1, 1).
		WillReturnRows(chatRows)

	messageRows := sqlmock.NewRows([]string{"id", "chat_id", "text", "created_at"}).
		AddRow(1, 1, "Message 1", createdAt.Add(time.Minute)).
		AddRow(2, 1, "Message 2", createdAt.Add(2*time.Minute))

	// С лимитом
	mock.ExpectQuery(`SELECT * FROM "messages" WHERE "messages"."chat_id" = $1 ORDER BY created_at DESC LIMIT $2`).
		WithArgs(1, 5).
		WillReturnRows(messageRows)

	chat, err := repo.GetByID(1, 5)

	assert.NoError(t, err)
	assert.NotNil(t, chat)
	assert.Equal(t, 1, chat.ID)
	assert.Len(t, chat.Messages, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	// УБРАЛИ условие "deleted_at" IS NULL
	mock.ExpectQuery(`SELECT * FROM "chats" WHERE "chats"."id" = $1 ORDER BY "chats"."id" LIMIT $2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at"}))

	chat, err := repo.GetByID(999, 20)

	assert.NoError(t, err)
	assert.Nil(t, chat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetByID_DatabaseError(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	mock.ExpectQuery(`SELECT * FROM "chats" WHERE "chats"."id" = $1 ORDER BY "chats"."id" LIMIT $2`).
		WithArgs(1, 1).
		WillReturnError(assert.AnError)

	chat, err := repo.GetByID(1, 20)

	assert.Error(t, err)
	assert.Nil(t, chat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_Delete_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "chats" WHERE "chats"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_Delete_Error(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewChatRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "chats" WHERE "chats"."id" = $1`).
		WithArgs(1).
		WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Delete(1)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
