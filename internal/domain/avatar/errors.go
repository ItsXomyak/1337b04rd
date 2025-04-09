package avatar

import "errors"

var (
	ErrFailedToFetchAvatar   = errors.New("failed to fetch avatar from external API")
	ErrNoAvailableAvatars    = errors.New("no available avatars left")
	ErrInvalidAvatarResponse = errors.New("invalid response format from avatar API")
	ErrAvatarAlreadyAssigned = errors.New("avatar already assigned to this session")
)
