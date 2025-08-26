package post_refresh

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
)

type UserService interface {
	RefreshTokens(ctx context.Context, refreshToken string) (*auth.Tokens, error)
}
