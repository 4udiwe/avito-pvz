package reception

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/reception_repo_mock.go

type ReceptionRepository interface {
	Open(ctx context.Context, pointID uuid.UUID) (entity.Reception, error)
	GetLastReceptionStatus(ctx context.Context, pointID uuid.UUID) (entity.ReceptionStatus, error)
	GetLastReceptionProductsAmount(ctx context.Context, pointID uuid.UUID) (int, error)
	CloseLastReception(ctx context.Context, pointID uuid.UUID) error
	CheckIfPointExists(ctx context.Context, pointID uuid.UUID) (bool, error)
}

type Metrics interface {
	Inc()
	ErrInc()
}
