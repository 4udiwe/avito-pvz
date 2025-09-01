package post_reception

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type ReceptionService interface {
	OpenReception(ctx context.Context, pointID uuid.UUID) (entity.Reception, error)
}
