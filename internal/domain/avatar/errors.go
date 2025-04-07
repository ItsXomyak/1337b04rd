package avatar

import "errors"

var (
	ErrAvatarNotFound = errors.New("avatar not found")
	ErrAvatarCreationFailed = errors.New("failed to create avatar")
	ErrAvatarUpdateFailed = errors.New("failed to update avatar")
	ErrAvatarDeleteFailed = errors.New("failed to delete avatar")
	ErrAvatarFileNotProvided = errors.New("avatar file not provided")
	ErrAvatarInvalidType = errors.New("invalid avatar type")
	ErrAvatarStorageError = errors.New("error occurred while storing avatar")
	ErrAvatarProcessingError = errors.New("error occurred while processing avatar")
)