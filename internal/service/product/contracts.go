package product

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	repo_product "github.com/4udiwe/avito-pvz/internal/repository/product"
	"github.com/google/uuid"
)

type ProductsRepository interface {
	Create(ctx context.Context, p repo_product.CreateProduct) (entity.Product, error)
	DeleteLastFromReception(ctx context.Context, pointID uuid.UUID) error
}

type ReceptionRepository interface {
	GetLastReceptionStatus(ctx context.Context, pointID uuid.UUID) (entity.ReceptionStatus, error)
}
