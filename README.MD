# Simple Chat API

Простой REST API для управления чатами и сообщениями, построенный на Go с использованием PostgreSQL и Docker.

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose

- Go 1.25+ (только для локальной разработки)

### Запуск проекта

#### 1. Клонируйте репозиторий:

```bash
git clone github.com/xGerfy/simple_api_chat
cd simpleApiChat
```

#### 2. Запустите все сервисы:

```bash
docker-compose up --build
```

#### 3. API будет доступно по адресу: http://localhost:8080

## Примеры использования

### Создание чата

```bash
curl -X POST http://localhost:8080/chats/ \
  -H "Content-Type: application/json" \
  -d '{"title": "Мой первый чат"}'
```

### Отправка сообщения

```bash
curl -X POST http://localhost:8080/chats/1/messages/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Привет, мир!"}'
```

### Получение чата

```bash
curl http://localhost:8080/chats/1?limit=10
```

### Удаление чата

```bash
curl -X DELETE http://localhost:8080/chats/1
```

## Структура проекта

```bash
simpleApiChat/
├── cmd/
│    └── server/
│    └── main.go         # Точка входа приложения
├── internal/
│    ├── app/            # Инициализация приложения
│    ├── config/         # Конфигурация
│    ├── handlers/       # HTTP обработчики
│    ├── models/         # Модели данных
│    ├── repository/     # Слой работы с БД
│    └── service/        # Бизнес-логика
├── migrations/          # SQL миграции для базы данных
├── tests/               # Тесты
├── Dockerfile           # Конфигурация Docker
├── docker-compose.yml   # Docker Compose конфигурация
├── go.mod               # Зависимости Go
└── README.md            # Этот файл
```

## Docker сервисы

| Сервис     | Порт | Описание               |
| ---------- | ---- | ---------------------- |
| postgres   | 5432 | PostgreSQL база данных |
| app        | 8080 | Go приложение с API    |
| migrations | -    | Выполнение миграций БД |

## API Endpoints

### 1. Создание чата

```text
POST /chats/
Content-Type: application/json

{
  "title": "Название чата"
}
```

#### Ограничения:

- title обязателен
- Длина: 1-200 символов
- Пробелы по краям автоматически обрезаются

### 2. Отправка сообщения

```text
POST /chats/{id}/messages/
Content-Type: application/json

{
  "text": "Текст сообщения"
}
```

#### Ограничения:

- text обязателен

- Длина: 1-5000 символов

- Пробелы по краям автоматически обрезаются

- Чат должен существовать (иначе 404)

### 3. Получение чата с сообщениями

```text
GET /chats/{id}?limit=20
```

#### Параметры:

- limit (опционально): количество сообщений (по умолчанию 20, максимум 100, если указать больше 100, все равно будет 100)

### 4. Удаление чата

```text
DELETE /chats/{id}
```

#### Примечание:

- Удаляет чат и все связанные сообщения (каскадное удаление).

## Модели данных

### Chat (чат)

```sql
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Message (Сообщение)

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);
```

## Команды разработки

### Docker команды

```bash
# Сборка и запуск
docker-compose up --build

# Остановка всех контейнеров
docker-compose down

# Остановка с удалением volumes (очистка БД)
docker-compose down -v

# Просмотр логов
docker-compose logs -f app
docker-compose logs -f postgres

# Выполнить команду в контейнере
docker-compose exec app sh
```

## Миграции

Для миграций в проекте используется Goose

### Up миграция

Происходит автоматически при запуске docker compose up

### Down миграция

```bash
goose -dir migrations postgres "host=localhost user=postgres password=postgres dbname=chatdb sslmode=disable" down
```

## Технологии

- Go 1.25 - основной язык программирования

- PostgreSQL 17 - реляционная база данных

- GORM - ORM для работы с БД

- Goose - инструмент миграций

- Docker - контейнеризация

- net/http - стандартная HTTP библиотека Go

- Testify - фреймворк для тестирования

## Ограничения и валидация

1. Chat title:

- Обязательное поле

- Длина: 1-200 символов

- Автоматический trim пробелов

2. Message text:

- Обязательное поле

- Длина: 1-5000 символов

- Автоматический trim пробелов

3. Чат должен существовать для отправки сообщений

4. Каскадное удаление - при удалении чата удаляются все его сообщения
