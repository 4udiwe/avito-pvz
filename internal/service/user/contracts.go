package user

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mocks.go -package=mocks

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) (entity.User, error)
}

type Auth interface {
	GenerateTokens(user entity.User) (*auth.Tokens, error)
	ValidateAccessToken(tokenString string) (*auth.TokenClaims, error)
	ValidateRefreshToken(tokenString string) (string, error)
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}
