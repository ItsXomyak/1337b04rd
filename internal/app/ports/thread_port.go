package ports

import (
	uuidHelper "1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/thread"
)

type ThreadPort interface {
	CreateThread(t *thread.Thread) error
	GetThreadByID(id uuidHelper.UUID) (*thread.Thread, error)
	UpdateThread(t *thread.Thread) error
	ListActiveThreads() ([]*thread.Thread, error) 
	ListAllThreads() ([]*thread.Thread, error)    
}