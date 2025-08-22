package product

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

type ProductsRepository interface {
	Create(ctx context.Context, pointID uuid.UUID, productType entity.ProductType) (entity.Product, error)
	DeleteLastFromReception(ctx context.Context, pointID uuid.UUID) error
}

type ReceptionRepository interface {
	GetLastReceptionStatus(ctx context.Context, pointID uuid.UUID) (entity.ReceptionStatus, error)
}
