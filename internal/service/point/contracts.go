package point

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

type PointRepository interface {
	Create(ctx context.Context, city string) (entity.Point, error)
	GetAll(ctx context.Context) ([]entity.Point, error)
}
