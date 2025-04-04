package shared

import "mime/multipart"

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UploadRequest struct {
	File *multipart.FileHeader `json:"file"`
}

type SummarylistResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	// Thumbnail   string `json:"thumbnail"`	// optional maybe for future
	AudioURL string `json:"audio_url"`
	// Transcript string `json:"transcript"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LlamaRequest models the request payload sent to the Llama API
type LlamaRequest struct {
	Messages     []map[string]string      `json:"messages"`
	Functions    []map[string]interface{} `json:"functions"`
	Stream       bool                     `json:"stream"`
	FunctionCall string                   `json:"function_call"` // If you want explicit function calls
}
