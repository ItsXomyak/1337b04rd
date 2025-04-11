package postgres

import (
	"database/sql"

	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/errors"
	"1337b04rd/internal/domain/thread"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

// CreateThread создает новый тред в базе данных.
func (r *ThreadRepository) CreateThread(t *thread.Thread) error {
	query := `
		INSERT INTO threads (id, title, content, image_url, session_id, created_at, last_commented, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query,
		t.ID,
		t.Title,
		t.Content,
		t.ImageURL,
		t.SessionID,
		t.CreatedAt,
		t.LastCommented,
		t.IsDeleted,
	)
	return err
}

func (r *ThreadRepository) GetThreadByID(id uuidHelper.UUID) (*thread.Thread, error) {
	query := `
		SELECT id, title, content, image_url, session_id, created_at, last_commented, is_deleted
		FROM threads
		WHERE id = $1`
	t := &thread.Thread{}
	var imageURL sql.NullString
	var lastCommented sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Content,
		&imageURL,
		&t.SessionID,
		&t.CreatedAt,
		&lastCommented,
		&t.IsDeleted,
	)
	if err == sql.ErrNoRows {
		return nil, errors.ErrThreadNotFound
	}
	if err != nil {
		return nil, err
	}

	if imageURL.Valid {
		t.ImageURL = &imageURL.String
	}
	if lastCommented.Valid {
		t.LastCommented = &lastCommented.Time
	}

	return t, nil
}

// UpdateThread обновляет существующий тред.
func (r *ThreadRepository) UpdateThread(t *thread.Thread) error {
	query := `
		UPDATE threads
		SET title = $2, content = $3, image_url = $4, session_id = $5, created_at = $6, last_commented = $7, is_deleted = $8
		WHERE id = $1`
	_, err := r.db.Exec(query,
		t.ID,
		t.Title,
		t.Content,
		t.ImageURL,
		t.SessionID,
		t.CreatedAt,
		t.LastCommented,
		t.IsDeleted,
	)
	return err
}

// ListActiveThreads возвращает список активных (не удаленных) тредов.
func (r *ThreadRepository) ListActiveThreads() ([]*thread.Thread, error) {
	query := `
		SELECT id, title, content, image_url, session_id, created_at, last_commented, is_deleted
		FROM threads
		WHERE is_deleted = FALSE`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*thread.Thread
	for rows.Next() {
		t := &thread.Thread{}
		var imageURL sql.NullString // Исправлено с sql.NullTime на sql.NullString
		var lastCommented sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&imageURL,
			&t.SessionID,
			&t.CreatedAt,
			&lastCommented,
			&t.IsDeleted,
		)
		if err != nil {
			return nil, err
		}

		if imageURL.Valid {
			t.ImageURL = &imageURL.String // Теперь String доступен
		}
		if lastCommented.Valid {
			t.LastCommented = &lastCommented.Time
		}

		threads = append(threads, t)
	}

	return threads, rows.Err()
}

// ListAllThreads возвращает список всех тредов, включая удаленные.
func (r *ThreadRepository) ListAllThreads() ([]*thread.Thread, error) {
	query := `
		SELECT id, title, content, image_url, session_id, created_at, last_commented, is_deleted
		FROM threads`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*thread.Thread
	for rows.Next() {
		t := &thread.Thread{}
		var imageURL sql.NullString // Исправлено с sql.NullTime на sql.NullString
		var lastCommented sql.NullTime

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&imageURL,
			&t.SessionID,
			&t.CreatedAt,
			&lastCommented,
			&t.IsDeleted,
		)
		if err != nil {
			return nil, err
		}

		if imageURL.Valid {
			t.ImageURL = &imageURL.String // Теперь String доступен
		}
		if lastCommented.Valid {
			t.LastCommented = &lastCommented.Time
		}

		threads = append(threads, t)
	}

	return threads, rows.Err()
}
