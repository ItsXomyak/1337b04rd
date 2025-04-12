package unit

// import (
// 	"testing"
// 	"time"

// 	"1337b04rd/internal/app/services"
// 	"1337b04rd/internal/domain/avatar"
// 	"1337b04rd/internal/domain/errors"
// 	"1337b04rd/internal/domain/session"
// )

// type mockAvatarService struct{}

// func (m *mockAvatarService) GetRandomAvatar() (*avatar.Avatar, error) {
// 	return &avatar.Avatar{
// 		URL:         "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
// 		DisplayName: "Rick Sanchez",
// 	}, nil
// }

// type mockSessionRepo struct {
// 	sessions []*session.Session
// }

// func (m *mockSessionRepo) GetSessionByID(id string) (*session.Session, error) {
// 	for _, s := range m.sessions {
// 		if s.ID.String() == id {
// 			return s, nil
// 		}
// 	}
// 	return nil, errors.ErrSessionNotFound
// }

// func (m *mockSessionRepo) CreateSession(s *session.Session) error {
// 	m.sessions = append(m.sessions, s)
// 	return nil
// }

// func (m *mockSessionRepo) DeleteExpired() error {
// 	return nil
// }

// func (m *mockSessionRepo) ListActiveSessions() ([]*session.Session, error) {
// 	return m.sessions, nil
// }

// func (m *mockSessionRepo) UpdateDisplayName(id string, name string) error {
// 	return nil
// }

// // --- тест ---
// func TestCreateNewSession(t *testing.T) {
// 	repo := &mockSessionRepo{}
// 	avatarSvc := &mockAvatarService{}
// 	ttl := time.Hour * 24 * 7

// 	svc := services.NewSessionService(repo, avatarSvc, ttl)

// 	sess, err := svc.CreateNew()
// 	if err != nil {
// 		t.Fatalf("CreateNew() returned error: %v", err)
// 	}

// 	if sess == nil {
// 		t.Fatal("CreateNew() returned nil session")
// 	}

// 	if sess.DisplayName != "Rick Sanchez" {
// 		t.Errorf("unexpected display name: got %q", sess.DisplayName)
// 	}

// 	if sess.AvatarURL == "" {
// 		t.Error("avatar URL should not be empty")
// 	}

// 	if sess.IsExpired() {
// 		t.Error("session should not be expired")
// 	}
// }
