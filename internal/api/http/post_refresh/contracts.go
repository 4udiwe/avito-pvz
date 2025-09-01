package post_refresh

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type UserService interface {
	RefreshTokens(ctx context.Context, refreshToken string) (*auth.Tokens, error)
}
