package get_points

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/entity"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type PointService interface {
	GetAllPointsFullInfo(ctx context.Context) ([]entity.PointFullInfo, error)
}
