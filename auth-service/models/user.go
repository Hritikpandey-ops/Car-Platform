package models

type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	IsVerified        bool   `json:"is_verified"`
	VerificationToken string `json:"-"`
}
