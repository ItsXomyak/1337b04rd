package post

import "errors"

var (
	ErrPostNotFound = errors.New("post not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternalError = errors.New("internal error")
	ErrThreadNotFound = errors.New("thread not found")
	ErrUserNotFound = errors.New("user not found")
	ErrNotAuthor = errors.New("not the author of the post")
	ErrPostTooLong = errors.New("post content too long")
	ErrInvalidThreadID = errors.New("invalid thread ID")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrPostAlreadyDeleted = errors.New("post already deleted")
	ErrPostAlreadyReported = errors.New("post already reported")
	ErrInvalidReportType = errors.New("invalid report type")
	ErrInvalidReportID = errors.New("invalid report ID")
	ErrReportNotFound = errors.New("report not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrInvalidCommentID = errors.New("invalid comment ID")
	ErrCommentTooLong = errors.New("comment content too long")
)