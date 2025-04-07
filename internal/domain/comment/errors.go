package comment

import "errors" 

var (
    ErrCommentNotFound = errors.New("comment not found")
		ErrCommentAlreadyExists = errors.New("comment already exists")
		ErrCommentAuthorMismatch = errors.New("comment author mismatch")
		ErrCommentContentInvalid = errors.New("comment content is invalid")
		ErrCommentThreadMismatch = errors.New("comment thread mismatch")
		ErrCommentPostMismatch = errors.New("comment post mismatch")
		ErrCommentParentMismatch = errors.New("comment parent mismatch")
		ErrCommentReplyMismatch = errors.New("comment reply mismatch")
		ErrCommentRatingOutOfRange = errors.New("comment rating out of range")
		ErrCommentStatusInvalid = errors.New("comment status is invalid")
		ErrCommentTimestampInvalid = errors.New("comment timestamp is invalid")
		ErrCommentVotesInvalid = errors.New("comment votes are invalid")
		ErrCommentEditedTimestampInvalid = errors.New("comment edited timestamp is invalid")
		ErrCommentEditedVotesInvalid = errors.New("comment edited votes are invalid")
		ErrCommentEditedStatusInvalid = errors.New("comment edited status is invalid")
		ErrCommentEditedTimestampOutOfRange = errors.New("comment edited timestamp out of range")

)