package model

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt uint   `json:"createdAt"`
	UpdatedAt uint   `json:"updatedAt"`
}
