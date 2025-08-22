package point

import (
	"context"
	"errors"

	"github.com/4udiwe/avito-pvz/internal/entity"
	"github.com/4udiwe/avito-pvz/internal/repository"
	"github.com/4udiwe/avito-pvz/pkg/transactor"
)

type Service struct {
	pointRepository PointRepository
	txManager       transactor.Transactor
}

func New(r PointRepository, tx transactor.Transactor) *Service {
	return &Service{
		pointRepository: r,
		txManager:       tx,
	}
}

func (s *Service) CreatePoint(ctx context.Context, city string) (entity.Point, error) {
	point, err := s.pointRepository.Create(ctx, city)

	if err != nil {
		if errors.Is(err, repository.ErrNoCityFound) {
			return entity.Point{}, ErrNoCityFound
		}
		return entity.Point{}, err
	}

	return point, nil
}

func (s *Service) GetAllPoints(ctx context.Context) ([]entity.Point, error) {
	points, err := s.pointRepository.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return points, nil
}


