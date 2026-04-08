-- +goose Up
CREATE TABLE IF NOT EXISTS submissions (
    id           BIGSERIAL PRIMARY KEY,
    task_id      BIGINT           NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id      BIGINT           NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code         TEXT             NOT NULL,
    status       VARCHAR(20)      NOT NULL DEFAULT 'pending',
    result       TEXT             NOT NULL DEFAULT '',
    score        INT              NOT NULL DEFAULT 0,
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_submissions_task_id ON submissions (task_id);
CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions (user_id);
