package get_points

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

type PointService interface {
	GetAllPointsFullInfo(ctx context.Context) ([]entity.PointFullInfo, error)
}
