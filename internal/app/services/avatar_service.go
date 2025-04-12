package services

import (
	"math/rand"
	"sync"
	"time"

	"1337b04rd/internal/adapters/rickmorty"
	"1337b04rd/internal/app/common/logger"
	"1337b04rd/internal/domain/avatar"
	"1337b04rd/internal/domain/errors"
)

type AvatarService struct {
	client   *rickmorty.Client
	shuffled []int
	index    int
	mu       sync.Mutex
}

func NewAvatarService(client *rickmorty.Client) *AvatarService {
	s := &AvatarService{
		client: client,
	}
	s.shuffleIDs()
	return s
}

func (s *AvatarService) GetRandomAvatar() (*avatar.Avatar, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.index >= len(s.shuffled) {
		s.shuffleIDs()
	}

	for attempts := 0; attempts < 10 && s.index < len(s.shuffled); attempts++ {
		id := s.shuffled[s.index]
		s.index++

		char, err := s.client.FetchCharacterByID(id)
		if err != nil {
			logger.Warn("failed to fetch character", "id", id, "error", err)
			continue
		}
		if char.Image == "" || char.Name == "" {
			logger.Warn("character missing data", "id", id, "image", char.Image, "name", char.Name)
			continue
		}

		return &avatar.Avatar{
			URL:         char.Image,
			DisplayName: char.Name,
		}, nil
	}

	logger.Error("GetRandomAvatar failed after multiple attempts")
	return nil, errors.ErrAfterMultipleAttempts
}

func (s *AvatarService) shuffleIDs() {
	const total = 826

	s.shuffled = make([]int, total)
	for i := 0; i < total; i++ {
		s.shuffled[i] = i + 1
	}

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	r.Shuffle(total, func(i, j int) {
		s.shuffled[i], s.shuffled[j] = s.shuffled[j], s.shuffled[i]
	})

	s.index = 0
}
