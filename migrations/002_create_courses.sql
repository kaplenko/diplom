-- +goose Up
CREATE TABLE IF NOT EXISTS courses (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(200)  NOT NULL,
    description TEXT          NOT NULL,
    created_by  BIGINT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
