package ports

import (
	"context"

	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/thread"
)

type ThreadPort interface {
	CreateThread(ctx context.Context,  t *thread.Thread) error
	GetThreadByID(ctx context.Context,  id uuidHelper.UUID) (*thread.Thread, error)
	UpdateThread(ctx context.Context,  t *thread.Thread) error
	ListActiveThreads(ctx context.Context) ([]*thread.Thread, error) 
	ListAllThreads(ctx context.Context) ([]*thread.Thread, error)    
}

