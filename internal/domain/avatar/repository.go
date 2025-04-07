package avatar

type AvatarRepository interface {
    Create(avatar Avatar) (int, error)
    GetByID(id int) (Avatar, error)
    Update(avatar Avatar) error
    Delete(id int) error
		ListForUser(userID int) ([]Avatar, error)
}
