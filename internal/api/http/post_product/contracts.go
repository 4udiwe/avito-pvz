package post_product

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

type ProductService interface {
	AddProduct(
		ctx context.Context,
		pointID uuid.UUID,
		productType entity.ProductType,
	) (entity.Product, error)
}
