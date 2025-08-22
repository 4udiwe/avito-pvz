package get_points

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

type PointService interface {
	GetAllPoints(ctx context.Context) ([]entity.Point, error)
}
