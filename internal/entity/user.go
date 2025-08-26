package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleModerator UserRole = "moderator"
	RoleEmployee  UserRole = "employee"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
