package patch_reception

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type ReceptionService interface {
	CloseReception(ctx context.Context, pointID uuid.UUID) error
}
