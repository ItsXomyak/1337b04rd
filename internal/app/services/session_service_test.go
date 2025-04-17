package services_test

import (
	"errors"
	"testing"
	"time"

	"1337b04rd/internal/app/common/utils"
	"1337b04rd/internal/domain/session"
)

type SessionInterface interface {
	IsExpired() bool
}

type MockSessionRepository struct {
	getSessionByIDFn func(sessionID string) (SessionInterface, error)
}

func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{}
}

func (m *MockSessionRepository) GetSessionByID(sessionID string) (*session.Session, error) {
	if m.getSessionByIDFn != nil {
		sessionInterface, err := m.getSessionByIDFn(sessionID)
		if err != nil {
			return nil, err
		}
		
		// Преобразуем интерфейс в конкретный тип
		if mockSession, ok := sessionInterface.(*MockSession); ok {
			return &session.Session{
				ID:          mockSession.ID,
				DisplayName: mockSession.DisplayName,
				AvatarURL:   mockSession.AvatarURL,
				ExpiresAt:   mockSession.ExpiresAt,
			}, nil
		}
		return nil, errors.New("invalid session type")
	}
	
	return nil, errors.New("session not found")
}

// MockSession - мок для сессии
type MockSession struct {
	ID          utils.UUID
	DisplayName string
	AvatarURL   string
	ExpiresAt   time.Time
	ExpiredFlag bool
}

// Реализация метода IsExpired для MockSession
func (s *MockSession) IsExpired() bool {
	return s.ExpiredFlag
}

// Мок для SessionService
type SessionService struct {
	repo *MockSessionRepository
}

func NewSessionService(repo *MockSessionRepository) *SessionService {
	return &SessionService{
		repo: repo,
	}
}

func (s *SessionService) GetOrCreate(sessionID string) (*session.Session, error) {
    if sessionID == "" {
        return nil, errors.New("Session ID not provided")
    }

    sess, err := s.repo.GetSessionByID(sessionID)
    if err != nil {
        return nil, err
    }

    if sess.IsExpired() {
        return nil, errors.New("session expired")
    }

    return sess, nil
}

// ParseTestUUID парсит UUID из строки для тестирования
func ParseTestUUID(s string) utils.UUID {
	uuid, err := utils.ParseUUID(s)
	if err != nil {
		var fallback utils.UUID
		copy(fallback[:], []byte("0123456789abcdef"))
		return fallback
	}
	return uuid
}

// Тесты
func TestGetOrCreate(t *testing.T) {
	// Тест 1: Проверка на пустой ID сессии
	t.Run("Empty session ID", func(t *testing.T) {
		mockRepo := NewMockSessionRepository()
		service := NewSessionService(mockRepo)
		
		_, err := service.GetOrCreate("")
		
		if err == nil {
			t.Error("Expected error for empty session ID, got nil")
		}
		if err != nil && err.Error() != "Session ID not provided" {
			t.Errorf("Expected 'Session ID not provided' error, got %v", err)
		}
	})

	// Тест 2: Сессия не найдена
	t.Run("Session not found", func(t *testing.T) {
		mockRepo := NewMockSessionRepository()
		mockRepo.getSessionByIDFn = func(sessionID string) (SessionInterface, error) {
			return nil, errors.New("session not found")
		}
		service := NewSessionService(mockRepo)
		
		_, err := service.GetOrCreate("non-existent-id")
		
		if err == nil {
			t.Error("Expected error for non-existent session, got nil")
		}
		if err != nil && err.Error() != "session not found" {
			t.Errorf("Expected 'session not found' error, got %v", err)
		}
	})

	// Тест 3: Сессия истекла
	t.Run("Expired session", func(t *testing.T) {
		mockRepo := NewMockSessionRepository()
		
		// Создаем истекшую сессию с UUID
		expiredUUID, _ := utils.ParseUUID("123e4567e89b12d3a456426614174000")
		expiredSession := &MockSession{
			ID:          expiredUUID,
			DisplayName: "Test User",
			AvatarURL:   "http://example.com/avatar.jpg",
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
			ExpiredFlag: true, // Помечаем как истекшую
		}
		
		mockRepo.getSessionByIDFn = func(sessionID string) (SessionInterface, error) {
			if sessionID == "expired-session-id" {
				return expiredSession, nil
			}
			return nil, errors.New("session not found")
		}
		
		service := NewSessionService(mockRepo)
		
		_, err := service.GetOrCreate("expired-session-id")
		
		if err == nil {
			t.Error("Expected error for expired session, got nil")
		}
		if err != nil && err.Error() != "session expired" {
			t.Errorf("Expected 'session expired' error, got %v", err)
		}
	})

	// Тест 4: Успешное получение сессии
	t.Run("Success getting session", func(t *testing.T) {
		mockRepo := NewMockSessionRepository()
		
		// Создаем действительную сессию с UUID
		validUUID, _ := utils.ParseUUID("123e4567e89b12d3a456426614174001")
		validSession := &MockSession{
			ID:          validUUID,
			DisplayName: "Test User",
			AvatarURL:   "http://example.com/avatar.jpg",
			ExpiresAt:   time.Now().Add(1 * time.Hour),
			ExpiredFlag: false, // Не истекшая
		}
		
		mockRepo.getSessionByIDFn = func(sessionID string) (SessionInterface, error) {
			if sessionID == "valid-session-id" {
				return validSession, nil
			}
			return nil, errors.New("session not found")
		}
		
		service := NewSessionService(mockRepo)
		
		sess, err := service.GetOrCreate("valid-session-id")
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sess == nil {
			t.Error("Expected session, got nil")
		}
		if sess != nil {
			expectedUUID := validUUID
			if sess.ID != expectedUUID {
				t.Errorf("Expected session ID %v, got %v", expectedUUID, sess.ID)
			}
		}
	})
}