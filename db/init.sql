-- init.sql

-- Clean up the database
DROP TABLE IF EXISTS posts;
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

-- posts
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    thread_id UUID NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    parent_post_id UUID REFERENCES posts(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT check_post_content_not_empty CHECK (char_length(content) > 0)
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
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION update_last_commented();

-- idxs
CREATE INDEX idx_posts_thread_id ON posts(thread_id);
CREATE INDEX idx_posts_parent_post_id ON posts(parent_post_id);
CREATE INDEX idx_threads_last_commented ON threads(last_commented);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
