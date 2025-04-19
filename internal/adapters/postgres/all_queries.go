package postgres

// thread repo
const (
	GetThreadByID = `
		SELECT id, title, content, image_url, session_id, 
		       created_at, last_commented, is_deleted
		FROM threads
		WHERE id = $1`

	CreateThread = `
		INSERT INTO threads (
			id, title, content, image_url, session_id, 
			created_at, last_commented, is_deleted
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	UpdateThread = `
		UPDATE threads
		SET title = $2, content = $3, image_url = $4, session_id = $5, 
		    created_at = $6, last_commented = $7, is_deleted = $8
		WHERE id = $1`

	ListActiveThreads = `
		SELECT id, title, content, image_url, session_id, 
		       created_at, last_commented, is_deleted
		FROM threads
		WHERE is_deleted = FALSE`

	ListAllThreads = `
		SELECT id, title, content, image_url, session_id, 
		       created_at, last_commented, is_deleted
		FROM threads`
)

// comment repo
const (
	CreateComment = `
		INSERT INTO comments (id, thread_id, parent_comment_id, content, image_url, session_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	GetCommentsByThreadID = `
		SELECT id, thread_id, parent_comment_id, content, image_url, session_id, created_at
		FROM comments
		WHERE thread_id = $1`
)

// session repo
const (
	CreateSession = `
		INSERT INTO sessions (id, avatar_url, display_name, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)`

	GetSessionByID = `
		SELECT id, avatar_url, display_name, created_at, expires_at
		FROM sessions
		WHERE id = $1`

	DeleteExpired = `
		DELETE FROM sessions
		WHERE expires_at < $1`

	ListActiveSessions = `
		SELECT id, avatar_url, display_name, created_at, expires_at
		FROM sessions`

	UpdateDisplayName = `UPDATE sessions SET display_name = $1 WHERE id = $2`
)
