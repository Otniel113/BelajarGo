package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Status    string    `json:"status" db:"status"` // "admin" or "member"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AuthRequest struct {
	Identity string `json:"identity"` // username or email
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
