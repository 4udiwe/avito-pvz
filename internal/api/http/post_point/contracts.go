package post_point

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

type PointService interface {
	CreatePoint(ctx context.Context, city string) (entity.Point, error)
}
