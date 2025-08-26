package post_dummy_login

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

type UserService interface {
	DummyLogin(ctx context.Context, role entity.UserRole) (string, error)
}
