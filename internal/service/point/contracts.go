package point

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/google/uuid"
)

type PointRepository interface {
	Create(ctx context.Context, city string) (entity.Point, error)
	GetAll(ctx context.Context) ([]entity.Point, error)
}

type ReceptionRepository interface {
	GetAllByPoint(ctx context.Context, pointID uuid.UUID) ([]entity.Reception, error)
}

type ProductRepository interface {
	GetAllByReception(ctx context.Context, receptionID uuid.UUID) ([]entity.Product, error)
}
