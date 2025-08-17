package entity

import "github.com/google/uuid"

type UserRole string

const (
	RoleModerator UserRole = "moderator"
	RoleEmployee  UserRole = "employee"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     UserRole  `db:"role"`
}
