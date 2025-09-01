package post_product

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type ProductService interface {
	AddProduct(
		ctx context.Context,
		pointID uuid.UUID,
		productType entity.ProductType,
	) (entity.Product, error)
}
