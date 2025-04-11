package unit_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"1337b04rd/internal/domain/thread"
)

func TestThreadService_ListActiveThreads(t *testing.T) {
    mockRepo := new(MockThreadPort)
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))
    svc := NewThreadService(mockRepo, logger)

    ctx := context.Background()
    mockRepo.On("ListActiveThreads", ctx).Return([]*thread.Thread{{ID: uuidHelper.NewUUID()}}, nil)

    threads, err := svc.ListActiveThreads(ctx)
    assert.NoError(t, err)
    assert.Len(t, threads, 1)
    mockRepo.AssertExpectations(t)
}