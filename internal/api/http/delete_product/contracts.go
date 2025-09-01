package delete_product

import (
	"context"

	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type ProductService interface {
	DeleteLastProductFromReception(ctx context.Context, pointID uuid.UUID) error
}
