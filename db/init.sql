-- Clean up the database
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS sessions;

-- sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    avatar_url TEXT NOT NULL,
    display_name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- threads
CREATE TABLE threads (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_commented TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT check_title_not_empty CHECK (char_length(title) > 0),
    CONSTRAINT check_content_not_empty CHECK (char_length(content) > 0)
);

-- comments
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    thread_id UUID NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT check_comment_content_not_empty CHECK (char_length(content) > 0)
);

-- triggers
CREATE OR REPLACE FUNCTION update_last_commented()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE threads
  SET last_commented = NEW.created_at
  WHERE id = NEW.thread_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_last_commented
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION update_last_commented();

-- indexes
CREATE INDEX idx_comments_thread_id ON comments(thread_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);
CREATE INDEX idx_threads_last_commented ON threads(last_commented);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
