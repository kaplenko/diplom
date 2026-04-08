-- +goose Up
CREATE TABLE IF NOT EXISTS progress (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id    BIGINT  NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lesson_id    BIGINT  NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    completed    BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (user_id, lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_progress_user_course ON progress (user_id, course_id);
