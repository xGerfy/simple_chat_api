-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    messages (
        id SERIAL PRIMARY KEY,
        chat_id INTEGER NOT NULL,
        text TEXT NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT fk_chat FOREIGN KEY (chat_id) REFERENCES chats (id) ON DELETE CASCADE
    );

CREATE INDEX idx_messages_chat_id ON messages (chat_id);

CREATE INDEX idx_messages_created_at ON messages (created_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;

-- +goose StatementEnd