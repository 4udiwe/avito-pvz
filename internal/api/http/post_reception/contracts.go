package post_reception

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

type ReceptionService interface {
	OpenReception(ctx context.Context, pointID uuid.UUID) (entity.Reception, error)
}
