-- +goose Up
CREATE TABLE IF NOT EXISTS lessons (
    id          BIGSERIAL PRIMARY KEY,
    course_id   BIGINT       NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title       VARCHAR(200) NOT NULL,
    content     TEXT         NOT NULL,
    order_index INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons (course_id);
