package post_point

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type PointService interface {
	CreatePoint(ctx context.Context, city string) (entity.Point, error)
}
