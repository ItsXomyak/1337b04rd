package postgres

import (
	"context"
	"database/sql"

	uuidHelper "1337b04rd/internal/app/common/utils"
	. "1337b04rd/internal/domain/errors"
	"1337b04rd/internal/domain/thread"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) *ThreadRepository {
	return &ThreadRepository{db: db}
}

func (r *ThreadRepository) CreateThread(ctx context.Context, t *thread.Thread) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, CreateThread,
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

func (r *ThreadRepository) GetThreadByID(ctx context.Context, id uuidHelper.UUID) (*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, GetThreadByID, id)
	t, err := scanThread(row)
	if err == sql.ErrNoRows {
		return nil, ErrThreadNotFound
	}
	return t, err
}

func (r *ThreadRepository) UpdateThread(ctx context.Context, t *thread.Thread) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, UpdateThread,
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

func (r *ThreadRepository) ListActiveThreads(ctx context.Context) ([]*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, ListActiveThreads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*thread.Thread
	for rows.Next() {
		t, err := scanThread(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, rows.Err()
}

func (r *ThreadRepository) ListAllThreads(ctx context.Context) ([]*thread.Thread, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, ListAllThreads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*thread.Thread
	for rows.Next() {
		t, err := scanThread(rows)
		if err != nil {
			return nil, err
		}
		threads = append(threads, t)
	}

	return threads, rows.Err()
}

func scanThread(scanner interface {
	Scan(dest ...interface{}) error
}) (*thread.Thread, error) {
	t := &thread.Thread{}
	var imageURL sql.NullString
	var lastCommented sql.NullTime

	err := scanner.Scan(
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
		t.ImageURL = &imageURL.String
	}
	if lastCommented.Valid {
		t.LastCommented = &lastCommented.Time
	}

	return t, nil
}