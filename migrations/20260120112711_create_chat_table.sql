-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    chats (
        id SERIAL PRIMARY KEY,
        title VARCHAR(200) NOT NULL,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE chats;

-- +goose StatementEnd