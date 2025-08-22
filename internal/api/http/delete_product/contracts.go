package delete_product

import (
	"context"

	"github.com/google/uuid"
)

type ProductService interface {
	DeleteLastProductFromReception(ctx context.Context, pointID uuid.UUID) error
}
