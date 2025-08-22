package patch_reception

import (
	"context"

	"github.com/google/uuid"
)

type ReceptionService interface {
	CloseReception(ctx context.Context, pointID uuid.UUID) error
}
