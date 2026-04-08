-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    id           BIGSERIAL PRIMARY KEY,
    lesson_id    BIGINT       NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    title        VARCHAR(200) NOT NULL,
    description  TEXT         NOT NULL,
    initial_code TEXT         NOT NULL DEFAULT '',
    test_cases   JSONB        NOT NULL DEFAULT '[]',
    difficulty   VARCHAR(20)  NOT NULL DEFAULT 'easy',
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_lesson_id ON tasks (lesson_id);
