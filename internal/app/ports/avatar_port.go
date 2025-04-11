package ports

import "1337b04rd/internal/domain/avatar"

type AvatarPort interface {
	GetRandomAvatar() (*avatar.Avatar, error)
}
