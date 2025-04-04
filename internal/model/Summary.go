package model

// keep the ID as uuid
type Recording struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	IsDeleted bool `json:"is_deleted"`
	CreatedAt uint `json:"created_at"`
	UpdatedAt uint `json:"updated_at"`
}
