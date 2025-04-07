package avatar

type Avatar struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	URL    string `json:"url"`
	Size   int    `json:"size"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version int    `json:"version"`
	Deleted bool   `json:"deleted"`
}